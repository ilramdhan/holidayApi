@echo off
echo ===============================================
echo 🔍 Testing Swagger Documentation
echo ===============================================

echo.
echo 1. Testing Swagger JSON endpoint...
curl -s http://localhost:8080/swagger/doc.json > nul
if %errorlevel% == 0 (
    echo ✅ Swagger JSON is accessible
) else (
    echo ❌ Swagger JSON is not accessible
    goto :end
)

echo.
echo 2. Testing Swagger UI endpoint...
curl -s http://localhost:8080/swagger/index.html > nul
if %errorlevel% == 0 (
    echo ✅ Swagger UI is accessible
) else (
    echo ❌ Swagger UI is not accessible
    goto :end
)

echo.
echo 3. Opening Swagger UI in browser...
start http://localhost:8080/swagger/index.html

echo.
echo 4. Alternative URLs to try:
echo    - http://localhost:8080/swagger/
echo    - http://localhost:8080/swagger/index.html
echo    - http://localhost:8080/swagger/doc.json

echo.
echo 5. If Swagger UI doesn't work:
echo    - Clear browser cache (Ctrl+F5)
echo    - Try incognito/private mode
echo    - Copy JSON from doc.json to https://editor.swagger.io/

:end
echo.
echo ===============================================
pause
