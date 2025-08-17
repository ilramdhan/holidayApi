@echo off
echo ===============================================
echo 📄 Exporting Swagger Documentation
echo ===============================================

echo.
echo Downloading Swagger JSON...
curl -s http://localhost:8080/swagger/doc.json -o swagger-export.json

if exist swagger-export.json (
    echo ✅ Swagger JSON exported to: swagger-export.json
    echo.
    echo You can now:
    echo 1. Open swagger-export.json in any text editor
    echo 2. Import it to https://editor.swagger.io/
    echo 3. Use it with other Swagger tools
    echo.
    echo File size: 
    for %%A in (swagger-export.json) do echo    %%~zA bytes
) else (
    echo ❌ Failed to export Swagger JSON
    echo Make sure the server is running on http://localhost:8080
)

echo.
echo ===============================================
pause
