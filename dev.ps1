# Turkcell Ev+Mobil Paket Danışmanı - Development Commands (PowerShell)
param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

function Show-Help {
    Write-Host "Available commands:" -ForegroundColor Green
    Write-Host "  .\dev.ps1 db-up      - Start Supabase local development"
    Write-Host "  .\dev.ps1 db-down    - Stop Supabase local development"
    Write-Host "  .\dev.ps1 api        - Start Go backend server"
    Write-Host "  .\dev.ps1 web        - Start Next.js frontend"
    Write-Host "  .\dev.ps1 test-api   - Run backend tests"
    Write-Host "  .\dev.ps1 test-web   - Run frontend tests"
    Write-Host "  .\dev.ps1 test       - Run all tests"
    Write-Host "  .\dev.ps1 seed       - Load sample data into database"
    Write-Host "  .\dev.ps1 clean      - Clean build artifacts"
}

switch ($Command) {
    "help" { Show-Help }
    "db-up" { 
        Set-Location db
        npx supabase start
        Set-Location ..
    }
    "db-down" { 
        Set-Location db
        npx supabase stop
        Set-Location ..
    }
    "api" { 
        Set-Location backend
        go run ./cmd/server
        Set-Location ..
    }
    "web" { 
        Set-Location frontend
        npm run dev
        Set-Location ..
    }
    "test-api" { 
        Set-Location backend
        go test ./...
        Set-Location ..
    }
    "test-web" { 
        Set-Location frontend
        npm test
        Set-Location ..
    }
    "test" { 
        Write-Host "Running backend tests..." -ForegroundColor Yellow
        Set-Location backend
        go test ./...
        Set-Location ..
        
        Write-Host "Running frontend tests..." -ForegroundColor Yellow
        Set-Location frontend
        npm test
        Set-Location ..
    }
    "seed" { 
        Set-Location db
        npx supabase db reset
        Set-Location ..
    }
    "clean" { 
        Write-Host "Cleaning build artifacts..." -ForegroundColor Yellow
        Set-Location backend
        go clean
        Set-Location ..
        
        Set-Location frontend
        if (Test-Path ".next") { Remove-Item -Recurse -Force ".next" }
        Set-Location ..
    }
    default { 
        Write-Host "Unknown command: $Command" -ForegroundColor Red
        Show-Help
    }
}
