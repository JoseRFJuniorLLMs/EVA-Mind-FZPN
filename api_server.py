#!/usr/bin/env python3
"""
EVA-Mind-FZPN - Python REST API Server
SPRINT 7 - Integration Layer
Integra com Go Integration Microservice
"""

from fastapi import FastAPI, Depends, HTTPException, status, Header
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import Optional, List
from datetime import datetime, timedelta
import psycopg2
from psycopg2.extras import RealDictCursor
import bcrypt
import jwt
import httpx
import os
from dotenv import load_dotenv

# Load environment
load_dotenv()

# Config
SECRET_KEY = os.getenv("JWT_SECRET_KEY", "your-secret-key-change-in-production")
ALGORITHM = "HS256"
ACCESS_TOKEN_EXPIRE_HOURS = 1

DB_CONFIG = {
    "host": os.getenv("DB_HOST", "104.248.219.200"),
    "port": int(os.getenv("DB_PORT", "5432")),
    "database": os.getenv("DB_NAME", "eva-db"),
    "user": os.getenv("DB_USER", "postgres"),
    "password": os.getenv("DB_PASSWORD", "Debian23@")
}

GO_SERVICE_URL = os.getenv("GO_SERVICE_URL", "http://localhost:8081")

# FastAPI app
app = FastAPI(
    title="EVA-Mind Integration API",
    description="REST API for EVA-Mind-FZPN (SPRINT 7)",
    version="1.0.0"
)

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# OAuth2
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="oauth/token")

# Pydantic Models
class TokenResponse(BaseModel):
    access_token: str
    token_type: str
    expires_in: int

class APIClient(BaseModel):
    id: str
    client_name: str
    client_type: str
    scopes: List[str]
    is_active: bool

class Patient(BaseModel):
    id: int
    name: str
    date_of_birth: str
    age: int
    gender: str
    email: Optional[str] = None
    phone: Optional[str] = None

class Assessment(BaseModel):
    id: str
    patient_id: int
    assessment_type: str
    total_score: Optional[int] = None
    severity: Optional[str] = None
    completed_at: datetime

# Database helpers
def get_db_connection():
    """Get database connection"""
    return psycopg2.connect(**DB_CONFIG, cursor_factory=RealDictCursor)

def create_access_token(data: dict):
    """Create JWT access token"""
    to_encode = data.copy()
    expire = datetime.utcnow() + timedelta(hours=ACCESS_TOKEN_EXPIRE_HOURS)
    to_encode.update({"exp": expire})
    encoded_jwt = jwt.encode(to_encode, SECRET_KEY, algorithm=ALGORITHM)
    return encoded_jwt

async def get_current_client(token: str = Depends(oauth2_scheme)) -> dict:
    """Validate JWT token and return client info"""
    credentials_exception = HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="Could not validate credentials",
        headers={"WWW-Authenticate": "Bearer"},
    )

    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        client_id: str = payload.get("client_id")
        if client_id is None:
            raise credentials_exception

        # Check if token is revoked
        conn = get_db_connection()
        cur = conn.cursor()

        cur.execute("""
            SELECT at.id, at.client_id, at.scopes, at.is_revoked,
                   ac.client_name, ac.is_active, ac.is_approved
            FROM api_tokens at
            JOIN api_clients ac ON ac.id = at.client_id
            WHERE at.access_token = %s
        """, (token,))

        token_data = cur.fetchone()
        cur.close()
        conn.close()

        if not token_data or token_data['is_revoked']:
            raise credentials_exception

        if not token_data['is_active'] or not token_data['is_approved']:
            raise HTTPException(
                status_code=status.HTTP_403_FORBIDDEN,
                detail="Client not active or not approved"
            )

        return dict(token_data)

    except jwt.ExpiredSignatureError:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Token expired"
        )
    except jwt.JWTError:
        raise credentials_exception

async def check_rate_limit(client_id: str):
    """Check if client exceeded rate limit"""
    conn = get_db_connection()
    cur = conn.cursor()

    # Count requests in last minute
    cur.execute("""
        SELECT COUNT(*) as count
        FROM api_request_logs
        WHERE client_id = %s
        AND timestamp > NOW() - INTERVAL '1 minute'
    """, (client_id,))

    count_last_minute = cur.fetchone()['count']

    # Get rate limit
    cur.execute("""
        SELECT rate_limit_per_minute
        FROM api_clients
        WHERE id = %s
    """, (client_id,))

    client = cur.fetchone()
    cur.close()
    conn.close()

    if count_last_minute >= client['rate_limit_per_minute']:
        raise HTTPException(
            status_code=status.HTTP_429_TOO_MANY_REQUESTS,
            detail="Rate limit exceeded"
        )

async def log_api_request(client_id: str, method: str, endpoint: str, status_code: int, response_time_ms: int):
    """Log API request to audit table"""
    conn = get_db_connection()
    cur = conn.cursor()

    cur.execute("""
        INSERT INTO api_request_logs (client_id, http_method, endpoint, http_status_code, response_time_ms, timestamp)
        VALUES (%s, %s, %s, %s, %s, NOW())
    """, (client_id, method, endpoint, status_code, response_time_ms))

    conn.commit()
    cur.close()
    conn.close()

# ============================================================================
# AUTHENTICATION ENDPOINTS
# ============================================================================

@app.post("/oauth/token", response_model=TokenResponse, tags=["Authentication"])
async def login(form_data: OAuth2PasswordRequestForm = Depends()):
    """
    OAuth2 Client Credentials Flow

    - **client_id**: API client ID
    - **client_secret**: API client secret
    """
    conn = get_db_connection()
    cur = conn.cursor()

    # Get client
    cur.execute("""
        SELECT id, client_name, client_secret_hash, scopes, is_active, is_approved
        FROM api_clients
        WHERE client_id = %s
    """, (form_data.username,))  # OAuth2PasswordRequestForm usa 'username' para client_id

    client = cur.fetchone()

    if not client:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid client credentials"
        )

    # Verify password
    if not bcrypt.checkpw(form_data.password.encode(), client['client_secret_hash'].encode()):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid client credentials"
        )

    if not client['is_active'] or not client['is_approved']:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Client not active or not approved"
        )

    # Create JWT token
    access_token = create_access_token(
        data={
            "client_id": str(client['id']),
            "scopes": client['scopes']
        }
    )

    # Save token to database
    expires_at = datetime.utcnow() + timedelta(hours=ACCESS_TOKEN_EXPIRE_HOURS)

    cur.execute("""
        INSERT INTO api_tokens (client_id, access_token, scopes, expires_at, created_at)
        VALUES (%s, %s, %s, %s, NOW())
    """, (client['id'], access_token, client['scopes'], expires_at))

    conn.commit()
    cur.close()
    conn.close()

    return TokenResponse(
        access_token=access_token,
        token_type="Bearer",
        expires_in=ACCESS_TOKEN_EXPIRE_HOURS * 3600
    )

# ============================================================================
# PATIENT ENDPOINTS
# ============================================================================

@app.get("/api/v1/patients/{patient_id}", response_model=Patient, tags=["Patients"])
async def get_patient(
    patient_id: int,
    current_client: dict = Depends(get_current_client)
):
    """
    Get patient by ID

    Requires scope: read:patients
    """
    # Check rate limit
    await check_rate_limit(current_client['client_id'])

    start_time = datetime.now()

    # Check scope
    if 'read:patients' not in current_client['scopes']:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Insufficient permissions"
        )

    # Call Go microservice
    try:
        async with httpx.AsyncClient() as client:
            response = await client.get(f"{GO_SERVICE_URL}/serialize/patient/{patient_id}")

        if response.status_code == 404:
            raise HTTPException(status_code=404, detail="Patient not found")

        response.raise_for_status()

        # Log request
        response_time_ms = int((datetime.now() - start_time).total_seconds() * 1000)
        await log_api_request(
            current_client['client_id'],
            "GET",
            f"/api/v1/patients/{patient_id}",
            response.status_code,
            response_time_ms
        )

        return response.json()

    except httpx.HTTPError as e:
        raise HTTPException(status_code=500, detail=f"Go service error: {str(e)}")

@app.get("/api/v1/patients", tags=["Patients"])
async def list_patients(
    limit: int = 10,
    offset: int = 0,
    current_client: dict = Depends(get_current_client)
):
    """
    List patients (paginated)

    Requires scope: read:patients
    """
    await check_rate_limit(current_client['client_id'])

    if 'read:patients' not in current_client['scopes']:
        raise HTTPException(status_code=403, detail="Insufficient permissions")

    conn = get_db_connection()
    cur = conn.cursor()

    # Get patients
    cur.execute("""
        SELECT id, name,
               EXTRACT(YEAR FROM AGE(date_of_birth::date))::int as age,
               gender
        FROM patients
        ORDER BY created_at DESC
        LIMIT %s OFFSET %s
    """, (limit, offset))

    patients = cur.fetchall()

    # Get total count
    cur.execute("SELECT COUNT(*) as count FROM patients")
    total_count = cur.fetchone()['count']

    cur.close()
    conn.close()

    return {
        "data": patients,
        "page": offset // limit + 1,
        "page_size": limit,
        "total_count": total_count,
        "has_next": offset + limit < total_count
    }

# ============================================================================
# ASSESSMENT ENDPOINTS
# ============================================================================

@app.get("/api/v1/assessments/{assessment_id}", response_model=Assessment, tags=["Assessments"])
async def get_assessment(
    assessment_id: str,
    current_client: dict = Depends(get_current_client)
):
    """
    Get assessment by ID

    Requires scope: read:assessments
    """
    await check_rate_limit(current_client['client_id'])

    if 'read:assessments' not in current_client['scopes']:
        raise HTTPException(status_code=403, detail="Insufficient permissions")

    # Call Go microservice
    try:
        async with httpx.AsyncClient() as client:
            response = await client.get(f"{GO_SERVICE_URL}/serialize/assessment/{assessment_id}")

        if response.status_code == 404:
            raise HTTPException(status_code=404, detail="Assessment not found")

        response.raise_for_status()
        return response.json()

    except httpx.HTTPError as e:
        raise HTTPException(status_code=500, detail=f"Go service error: {str(e)}")

# ============================================================================
# FHIR ENDPOINTS
# ============================================================================

@app.get("/api/v1/fhir/patients/{patient_id}", tags=["FHIR"])
async def get_patient_fhir(
    patient_id: int,
    current_client: dict = Depends(get_current_client)
):
    """
    Get patient in FHIR R4 format

    Requires scope: read:patients or export:data
    """
    await check_rate_limit(current_client['client_id'])

    if 'read:patients' not in current_client['scopes'] and 'export:data' not in current_client['scopes']:
        raise HTTPException(status_code=403, detail="Insufficient permissions")

    # Call Go microservice
    try:
        async with httpx.AsyncClient() as client:
            response = await client.get(f"{GO_SERVICE_URL}/fhir/patient/{patient_id}")

        response.raise_for_status()
        return response.json()

    except httpx.HTTPError as e:
        raise HTTPException(status_code=500, detail=f"Go service error: {str(e)}")

@app.get("/api/v1/fhir/bundle/{patient_id}", tags=["FHIR"])
async def get_fhir_bundle(
    patient_id: int,
    current_client: dict = Depends(get_current_client)
):
    """
    Get FHIR Bundle for patient (Patient + Observations)

    Requires scope: export:data
    """
    await check_rate_limit(current_client['client_id'])

    if 'export:data' not in current_client['scopes']:
        raise HTTPException(status_code=403, detail="Insufficient permissions")

    # Call Go microservice
    try:
        async with httpx.AsyncClient() as client:
            response = await client.get(f"{GO_SERVICE_URL}/fhir/bundle/{patient_id}")

        response.raise_for_status()
        return response.json()

    except httpx.HTTPError as e:
        raise HTTPException(status_code=500, detail=f"Go service error: {str(e)}")

# ============================================================================
# EXPORT ENDPOINTS
# ============================================================================

@app.get("/api/v1/export/lgpd/{patient_id}", tags=["Export"])
async def export_lgpd(
    patient_id: int,
    current_client: dict = Depends(get_current_client)
):
    """
    Export patient data (LGPD/GDPR portability)

    Requires scope: export:data
    """
    await check_rate_limit(current_client['client_id'])

    if 'export:data' not in current_client['scopes']:
        raise HTTPException(status_code=403, detail="Insufficient permissions")

    # Call Go microservice
    try:
        async with httpx.AsyncClient() as client:
            response = await client.get(f"{GO_SERVICE_URL}/export/lgpd/{patient_id}")

        response.raise_for_status()
        return response.json()

    except httpx.HTTPError as e:
        raise HTTPException(status_code=500, detail=f"Go service error: {str(e)}")

# ============================================================================
# HEALTH & INFO
# ============================================================================

@app.get("/health", tags=["System"])
async def health_check():
    """Health check endpoint"""
    # Check Go service
    go_service_healthy = False
    try:
        async with httpx.AsyncClient() as client:
            response = await client.get(f"{GO_SERVICE_URL}/health", timeout=2.0)
            go_service_healthy = response.status_code == 200
    except:
        pass

    return {
        "status": "healthy" if go_service_healthy else "degraded",
        "api_server": "healthy",
        "go_service": "healthy" if go_service_healthy else "unhealthy",
        "timestamp": datetime.utcnow()
    }

@app.get("/", tags=["System"])
async def root():
    """API info"""
    return {
        "name": "EVA-Mind Integration API",
        "version": "1.0.0",
        "sprint": "SPRINT 7 - Integration Layer",
        "docs": "/docs",
        "health": "/health"
    }

# ============================================================================
# RUN SERVER
# ============================================================================

if __name__ == "__main__":
    import uvicorn

    print("=" * 60)
    print("🚀 EVA-Mind Integration API Server")
    print("=" * 60)
    print(f"API Server: http://localhost:8000")
    print(f"API Docs: http://localhost:8000/docs")
    print(f"Go Service: {GO_SERVICE_URL}")
    print(f"Database: {DB_CONFIG['host']}:{DB_CONFIG['port']}/{DB_CONFIG['database']}")
    print("=" * 60)

    uvicorn.run(app, host="0.0.0.0", port=8000)
