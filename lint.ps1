# PowerShell script for running linters

# Function to check if a command exists
function Test-Command($cmdname) {
    return [bool](Get-Command -Name $cmdname -ErrorAction SilentlyContinue)
}

# Function to clear linter caches
function Clear-LinterCaches {
    Write-Host "Clearing linter caches..." -ForegroundColor Yellow
    Remove-Item -Path "$env:LOCALAPPDATA\revive" -Recurse -Force -ErrorAction SilentlyContinue
    Remove-Item -Path "$env:LOCALAPPDATA\golangci-lint" -Recurse -Force -ErrorAction SilentlyContinue
    Write-Host "Caches cleared" -ForegroundColor Green
}

# Install tools if not present
function Install-Tools {
    Write-Host "Checking and installing required tools..." -ForegroundColor Yellow
    
    if (-not (Test-Command "golangci-lint")) {
        Write-Host "Installing golangci-lint..." -ForegroundColor Cyan
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    }
    
    if (-not (Test-Command "revive")) {
        Write-Host "Installing revive..." -ForegroundColor Cyan
        go install github.com/mgechev/revive@latest
    }

    # Verify installations
    Write-Host "`nVerifying installations:" -ForegroundColor Yellow
    Write-Host "golangci-lint version: " -NoNewline
    golangci-lint --version
    Write-Host "revive version: " -NoNewline
    revive --version
}

# Run all linters
function Run-Linters {
    Write-Host "Running linters..." -ForegroundColor Yellow
    $hasErrors = $false
    
    Write-Host "`nRunning golangci-lint..." -ForegroundColor Cyan
    golangci-lint run --config=.golangci.yml ./...
    if ($LASTEXITCODE -ne 0) {
        $hasErrors = $true
    }
    
    Write-Host "`nRunning revive..." -ForegroundColor Cyan
    revive -config revive.toml -formatter friendly ./...
    if ($LASTEXITCODE -ne 0) {
        $hasErrors = $true
    }

    if ($hasErrors) {
        Write-Host "`nSome linting checks failed. Try running option 4 to fix common issues." -ForegroundColor Yellow
    }
}

# Fix common issues
function Fix-LintIssues {
    Write-Host "Fixing common issues..." -ForegroundColor Yellow
    
    Write-Host "Running gofmt..." -ForegroundColor Cyan
    gofmt -w .
    
    Write-Host "Running golangci-lint fix..." -ForegroundColor Cyan
    try {
        golangci-lint run --fix --config=.golangci.yml ./... 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Host "Fixed common issues successfully" -ForegroundColor Green
        } else {
            Write-Host "Some issues could not be fixed automatically" -ForegroundColor Yellow
        }
    } catch {
        Write-Host "Error running golangci-lint fix: $_" -ForegroundColor Red
        Write-Host "Try running option 2 to clear caches and then try again" -ForegroundColor Yellow
    }
}

# Run tests
function Run-Tests {
    Write-Host "Running tests..." -ForegroundColor Yellow
    go test -v -race ./...
}

# Main menu
function Show-Menu {
    Write-Host "`nGo Development Tools" -ForegroundColor Green
    Write-Host "===================" -ForegroundColor Green
    Write-Host "1: Install/Update tools"
    Write-Host "2: Clear linter caches"
    Write-Host "3: Run all linters"
    Write-Host "4: Fix common lint issues"
    Write-Host "5: Run tests"
    Write-Host "6: Run revive only"
    Write-Host "Q: Quit"
    Write-Host "===================" -ForegroundColor Green
}

# Check for config files
if (-not (Test-Path ".golangci.yml")) {
    Write-Host "Warning: .golangci.yml not found in current directory" -ForegroundColor Red
}
if (-not (Test-Path "revive.toml")) {
    Write-Host "Warning: revive.toml not found in current directory" -ForegroundColor Red
}

# Main loop
do {
    Show-Menu
    $input = Read-Host "Please make a selection"
    switch ($input) {
        '1' { Install-Tools }
        '2' { Clear-LinterCaches }
        '3' { Run-Linters }
        '4' { Fix-LintIssues }
        '5' { Run-Tests }
        '6' { revive -config revive.toml -formatter friendly ./... }
        'q' { return }
    }
    if ($input -ne 'q') {
        Write-Host "`nPress any key to continue..."
        $null = $Host.UI.RawUI.ReadKey('NoEcho,IncludeKeyDown')
    }
} until ($input -eq 'q') 