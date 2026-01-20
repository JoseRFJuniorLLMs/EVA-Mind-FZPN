@echo off
REM Script de deploy rápido - Windows
echo 📦 Adicionando arquivos...
git add .

echo 💾 Commitando...
git commit -m "deploy"

echo 🚀 Enviando para o servidor...
git push

echo ✅ Pronto! Agora rode no servidor:
echo    cd ~/EVA-Mind-FZPN
echo    ./deploy.sh
