import os
import time
import json
import argparse
from datetime import datetime, timedelta
import google.generativeai as genai
from neo4j import GraphDatabase
from qdrant_client import QdrantClient
from qdrant_client.http import models
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Configuration
GOOGLE_API_KEY = os.getenv("GOOGLE_API_KEY")
QDRANT_HOST = os.getenv("QDRANT_HOST", "localhost")
QDRANT_PORT = int(os.getenv("QDRANT_PORT", "6333"))
NEO4J_URI = os.getenv("NEO4J_URI", "bolt://localhost:7687")
NEO4J_USER = os.getenv("NEO4J_USER", "neo4j")
NEO4J_PASSWORD = os.getenv("NEO4J_PASSWORD")

# Setup Gemini
genai.configure(api_key=GOOGLE_API_KEY)
model = genai.GenerativeModel('gemini-1.5-pro-latest') # Use a capable model for generation

# Setup Clients
qdrant = QdrantClient(host=QDRANT_HOST, port=QDRANT_PORT)
neo4j_driver = GraphDatabase.driver(NEO4J_URI, auth=(NEO4J_USER, NEO4J_PASSWORD))

# Bio Context - The Soul of the Simulation
BIO_CONTEXT = {
    "name": "Jose R F Junior",
    "cpf": "645.254.302-49",
    "family": {
        "daughters": [
            {"name": "Electra", "birth": "2008-05-15", "desc": "Living in Brazil, digital native, connects via video."},
            {"name": "Coraline", "birth": "2010-08-21", "desc": "Living in Brazil, deeply bonded, creative."},
            {"name": "Elizabeth", "birth": "2023-05-19", "desc": "The baby (Caçula). Born in 2023. Source of immense joy and renewal."}
        ]
    },
    "career": [
        {"period": "2000-2025", "role": "Técnico do Seguro Social (INSS)", "desc": "Matrícula 1634972. Heavy bureaucracy, fighting for efficiency."},
        {"period": "2026-Present", "role": "Coordenador de Estudos e Pesquisas (MinC)", "desc": "Nível FCE 1.10. Strategic management of cultural data. Transitioned via Ofício 107/2026."}
    ],
    "health": [
        {"year": "2016", "event": "Diagnosed with Hypertension (Primary)", "meds": ["Losartana 50mg"]},
        {"year": "2019", "event": "Fall in bathroom (3 AM)", "consequence": "Hip trauma, fear of loss of autonomy."}
    ],
    "literary_works": [
        "Sarmoung (Sci-Fi/Philosophy - Magnum Opus)",
        "O Arquivo Akáshico da Burocracia (Fiction)",
        "Gestão da Informação no Serviço Social: O Dado como Afeto (Technical)",
        "Manifesto da Invasão Antropomórfica (AI Ethics)"
    ]
}

def generate_yearly_summary(year, cpf, life_stage):
    # Calculate daughters' ages
    daughters_status = []
    for d in BIO_CONTEXT["family"]["daughters"]:
        b_year = int(d["birth"][:4])
        age = year - b_year
        if age >= 0:
            daughters_status.append(f"{d['name']} is {age} years old.")
        elif age == -1:
            daughters_status.append(f"{d['name']} is expected to be born next year.")
        elif age == 0:
            daughters_status.append(f"{d['name']} WAS BORN THIS YEAR ({d['birth']})! A major event.")

    # Determine Career
    career_role = "INSS Technician"
    if year >= 2026:
        career_role = "Coordinator at MinC (Ministry of Culture)"
    
    prompt = f"""
    Acting as a creative storyteller/biographer, generate a detailed annual summary for:
    Name: {BIO_CONTEXT['name']}
    Year: {year}
    Role: {career_role}
    
    HARD FACTS (Must Respect):
    - Family Status: {', '.join(daughters_status)}
    - Health History: {json.dumps(BIO_CONTEXT['health'])}
    - Literary Focus: Writing "{BIO_CONTEXT['literary_works']}"
    
    NARRATIVE ARC FOR {year}:
    {life_stage}
    
    Requirements:
    1. Narrative (1st Person Perspective preferred for 'memory'): "I felt...", "Electra called me...". Approx 150 words.
    2. Events: 3-5 distinct events (Medical, Family, Career).
    3. Emotion: Dominant feeling.
    
    Output JSON:
    {{
        "narrative": "...",
        "events": [
            {{"date": "{year}-MM-DD", "description": "...", "severity": "LOW|MEDIUM|HIGH"}}
        ],
        "dominant_emotion": "..."
    }}
    """
    
    response = model.generate_content(prompt, generation_config={"response_mime_type": "application/json"})
    return json.loads(response.text)

def embedding_model(text):
    result = genai.embed_content(
        model="models/text-embedding-004",
        content=text,
        task_type="retrieval_document",
        title="Patient History"
    )
    return result['embedding']

def seed_qdrant(data, collection_name="episodic_memory"):
    # Ensure collection exists
    try:
        qdrant.get_collection(collection_name)
    except:
        qdrant.create_collection(
            collection_name=collection_name,
            vectors_config=models.VectorParams(size=768, distance=models.Distance.COSINE)
        )
    
    # Create vector
    vector = embedding_model(data["narrative"])
    
    # Upsert
    point_id = int(time.time() * 1000) # Simple ID generation
    qdrant.upsert(
        collection_name=collection_name,
        points=[
            models.PointStruct(
                id=point_id,
                vector=vector,
                payload={
                    "year": data["year"],
                    "cpf": data["cpf"],
                    "type": "yearly_summary",
                    "content": data["narrative"],
                    "emotion": data["dominant_emotion"]
                }
            )
        ]
    )
    print(f"✅ [Qdrant] Seeded year {data['year']}")

def seed_neo4j(data):
    query = """
    MERGE (p:Person {cpf: $cpf})
    MERGE (y:Year {value: $year})
    MERGE (p)-[:LIVED_THROUGH]->(y)
    
    WITH p, y
    UNWIND $events AS event
    CREATE (e:Event {
        date: event.date,
        description: event.description,
        severity: event.severity
    })
    MERGE (p)-[:EXPERIENCED]->(e)
    MERGE (e)-[:HAPPENED_IN]->(y)
    """
    
    with neo4j_driver.session() as session:
        session.run(query, cpf=data["cpf"], year=data["year"], events=data["events"])
        
    print(f"✅ [Neo4j] Seeded events for {data['year']}")

def main():
    parser = argparse.ArgumentParser(description="Seed synthetic history for EVA-Mind")
    parser.add_argument("--cpf", required=True, help="Patient CPF")
    parser.add_argument("--start_year", type=int, default=2016, help="Start year of history")
    parser.add_argument("--years", type=int, default=10, help="Number of years to generate")
    
    args = parser.parse_args()
    
    print(f"🚀 Starting History Seeding for {args.cpf}...")
    
    for i in range(args.years):
        year = args.start_year + i
        
        # Determine Life Stage context
        if year < 2019:
            stage = "Pre-Diagnosis / Denial"
        elif 2019 <= year <= 2021:
            stage = "Crisis / Vulnerability / First EVA Contact"
        else:
            stage = "Stabilization / Acceptance / Friendship with EVA"
            
        print(f"Generating {year} ({stage})...")
        
        try:
            summary = generate_yearly_summary(year, args.cpf, stage)
            
            # Combine for DBs
            seed_data = {
                "cpf": args.cpf,
                "year": year,
                "narrative": summary["narrative"],
                "dominant_emotion": summary["dominant_emotion"],
                "events": summary["events"]
            }
            
            seed_qdrant(seed_data)
            seed_neo4j(seed_data)
            
            # Rate limit friendly
            time.sleep(2)
            
        except Exception as e:
            print(f"❌ Error generating {year}: {e}")

    neo4j_driver.close()
    print("✨ Seeding Complete!")

if __name__ == "__main__":
    main()
