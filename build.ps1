# PowerShell build script for gsecutil
# Builds binaries for multiple platforms
# Compatible with Makefile targets

param(
    [Parameter(Position=0)]
    [string]$Target = "default"
)

# Set error action preference
$ErrorActionPreference = "Stop"

# Constants
$BinaryName = "gsecutil"
$Version = if ($env:VERSION) {
    $env:VERSION
} elseif (Test-Path "VERSION") {
    Get-Content "VERSION" -Raw | ForEach-Object { $_.Trim() }
} else {
    "1.0.0"
}
$BuildDir = "build"

# Define build targets (OS/ARCH combinations)
$AllTargets = @(
    @{OS="linux"; ARCH="amd64"},
    @{OS="linux"; ARCH="arm64"},
    @{OS="darwin"; ARCH="amd64"},
    @{OS="darwin"; ARCH="arm64"},
    @{OS="windows"; ARCH="amd64"}
)

# Function to show help
function Show-Help {
    Write-Host "Usage: .\\build.ps1 [target]" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Available targets:" -ForegroundColor Blue
    Write-Host "  default       - Build for current platform (with debug info)" -ForegroundColor Green
    Write-Host "  all           - Build for all platforms (release builds)" -ForegroundColor Green
    Write-Host "  clean         - Clean build directory" -ForegroundColor Green
    Write-Host "  linux         - Build for Linux amd64 (release)" -ForegroundColor Green
    Write-Host "  linux-arm64   - Build for Linux arm64 (release)" -ForegroundColor Green
    Write-Host "  windows       - Build for Windows amd64 (release)" -ForegroundColor Green
    Write-Host "  darwin        - Build for macOS amd64 (release)" -ForegroundColor Green
    Write-Host "  darwin-arm64  - Build for macOS arm64 (release)" -ForegroundColor Green
    Write-Host "  test          - Run tests" -ForegroundColor Green
    Write-Host "  fmt           - Format code" -ForegroundColor Green
    Write-Host "  vet           - Run go vet" -ForegroundColor Green
    Write-Host "  deps          - Install dependencies" -ForegroundColor Green
    Write-Host "  install       - Install locally" -ForegroundColor Green
    Write-Host "  dev           - Development build and show help" -ForegroundColor Green
    Write-Host "  help          - Show this help message" -ForegroundColor Green
}

# Function to clean build directory
function Clean-BuildDir {
    Write-Host "Cleaning build directory..." -ForegroundColor Blue

    if (Test-Path $BuildDir) {
        Remove-Item $BuildDir -Recurse -Force
    }

    & go clean
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to run go clean"
    }

    Write-Host "✓ Clean completed" -ForegroundColor Green
}

# Function to create build directory if it doesn't exist
function Create-BuildDir {
    if (-not (Test-Path $BuildDir)) {
        Write-Host "Creating build directory..." -ForegroundColor Blue
        New-Item -ItemType Directory -Path $BuildDir -Force | Out-Null
        Write-Host "✓ Build directory created" -ForegroundColor Green
    }
}

# Function to build for a specific OS/ARCH
function Build-ForTarget {
    param(
        [string]$GOOS,
        [string]$GOARCH,
        [string]$OutputSuffix = ""
    )

    Write-Host "Building for $GOOS/$GOARCH..." -ForegroundColor Yellow

    if ($OutputSuffix) {
        $OutputName = "$BinaryName$OutputSuffix"
    } else {
        $OutputName = "$BinaryName-$GOOS-$GOARCH"
    }

    if ($GOOS -eq "windows") {
        $OutputName = "$OutputName.exe"
    }

    $OutputPath = Join-Path $BuildDir $OutputName

    try {
        $env:GOOS = $GOOS
        $env:GOARCH = $GOARCH

        & go build -ldflags "-X main.Version=$Version -s -w" -o $OutputPath .

        if ($LASTEXITCODE -ne 0) {
            throw "Build failed"
        }

        Write-Host "✓ Built $OutputName" -ForegroundColor Green
        return $true
    }
    catch {
        Write-Host "✗ Failed to build $OutputName" -ForegroundColor Red
        return $false
    }
    finally {
        # Clean up environment variables
        Remove-Item env:GOOS -ErrorAction SilentlyContinue
        Remove-Item env:GOARCH -ErrorAction SilentlyContinue
    }
}

# Function to build for all platforms
function Build-All {
    Write-Host "Building $BinaryName v$Version for all platforms" -ForegroundColor Blue
    Create-BuildDir

    $Failed = $false

    # Build for each target
    foreach ($Target in $AllTargets) {
        $Success = Build-ForTarget -GOOS $Target.OS -GOARCH $Target.ARCH
        if (-not $Success) {
            $Failed = $true
        }
    }

    if (-not $Failed) {
        Show-BuildSummary
        Write-Host "All builds completed successfully!" -ForegroundColor Green
        return $true
    } else {
        Write-Host "Some builds failed. See above for details." -ForegroundColor Red
        return $false
    }
}

# Function to build for current platform (development build with debug info)
function Build-Current {
    Write-Host "Building $BinaryName v$Version for current platform (with debug info)" -ForegroundColor Blue
    Create-BuildDir

    $OutputPath = Join-Path $BuildDir "$BinaryName.exe"

    try {
        # Keep debug info for development builds (no -s -w flags)
        & go build -ldflags "-X main.Version=$Version" -o $OutputPath .

        if ($LASTEXITCODE -ne 0) {
            throw "Build failed"
        }

        Write-Host "✓ Built $BinaryName for current platform (development build)" -ForegroundColor Green
        Show-BuildSummary
        return $true
    }
    catch {
        Write-Host "✗ Failed to build $BinaryName for current platform" -ForegroundColor Red
        return $false
    }
}

# Function to run tests
function Run-Tests {
    Write-Host "Running tests..." -ForegroundColor Blue

    & go test ./...

    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Tests passed" -ForegroundColor Green
        return $true
    } else {
        Write-Host "✗ Tests failed" -ForegroundColor Red
        return $false
    }
}

# Function to format code
function Format-Code {
    Write-Host "Formatting code..." -ForegroundColor Blue

    & go fmt ./...

    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Code formatted" -ForegroundColor Green
        return $true
    } else {
        Write-Host "✗ Code formatting failed" -ForegroundColor Red
        return $false
    }
}

# Function to run vet
function Run-Vet {
    Write-Host "Running go vet..." -ForegroundColor Blue

    & go vet ./...

    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Go vet passed" -ForegroundColor Green
        return $true
    } else {
        Write-Host "✗ Go vet found issues" -ForegroundColor Red
        return $false
    }
}

# Function to install dependencies
function Install-Dependencies {
    Write-Host "Installing dependencies..." -ForegroundColor Blue

    & go mod tidy
    if ($LASTEXITCODE -ne 0) {
        Write-Host "✗ go mod tidy failed" -ForegroundColor Red
        return $false
    }

    & go mod download
    if ($LASTEXITCODE -ne 0) {
        Write-Host "✗ go mod download failed" -ForegroundColor Red
        return $false
    }

    Write-Host "✓ Dependencies installed" -ForegroundColor Green
    return $true
}

# Function to install locally
function Install-Locally {
    Write-Host "Installing $BinaryName v$Version locally..." -ForegroundColor Blue

    & go install -ldflags "-X main.Version=$Version" .

    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Installed successfully" -ForegroundColor Green
        return $true
    } else {
        Write-Host "✗ Installation failed" -ForegroundColor Red
        return $false
    }
}

# Function for development build and run
function Run-Dev {
    $Success = Build-Current
    if ($Success) {
        Write-Host "Running $BinaryName help..." -ForegroundColor Blue
        $BinaryPath = Join-Path $BuildDir "$BinaryName.exe"
        & $BinaryPath --help
    }
}

# Function to show build summary
function Show-BuildSummary {
    if (-not (Test-Path $BuildDir)) {
        return
    }

    Write-Host ""
    Write-Host "Build completed!" -ForegroundColor Blue
    Write-Host "Binaries are available in the $BuildDir directory:" -ForegroundColor Blue
    Write-Host ""

    Get-ChildItem $BuildDir | Format-Table Name, Length, LastWriteTime -AutoSize

    Write-Host ""
    Write-Host "Binary sizes:" -ForegroundColor Blue

    Get-ChildItem $BuildDir -File | ForEach-Object {
        $SizeMB = [math]::Round($_.Length / 1MB, 2)
        $SizeKB = [math]::Round($_.Length / 1KB, 0)

        if ($SizeMB -gt 1) {
            Write-Host "  $($_.Name): ${SizeMB} MB"
        } else {
            Write-Host "  $($_.Name): ${SizeKB} KB"
        }
    }
}

# Main logic - process command line argument
try {
    switch ($Target.ToLower()) {
        "default" {
            $Success = Build-Current
        }
        "all" {
            $Success = Build-All
        }
        "clean" {
            Clean-BuildDir
            $Success = $true
        }
        "linux" {
            Create-BuildDir
            $Success = Build-ForTarget -GOOS "linux" -GOARCH "amd64"
            if ($Success) { Show-BuildSummary }
        }
        "linux-arm64" {
            Create-BuildDir
            $Success = Build-ForTarget -GOOS "linux" -GOARCH "arm64"
            if ($Success) { Show-BuildSummary }
        }
        "windows" {
            Create-BuildDir
            $Success = Build-ForTarget -GOOS "windows" -GOARCH "amd64"
            if ($Success) { Show-BuildSummary }
        }
        "darwin" {
            Create-BuildDir
            $Success = Build-ForTarget -GOOS "darwin" -GOARCH "amd64"
            if ($Success) { Show-BuildSummary }
        }
        "darwin-arm64" {
            Create-BuildDir
            $Success = Build-ForTarget -GOOS "darwin" -GOARCH "arm64"
            if ($Success) { Show-BuildSummary }
        }
        "test" {
            $Success = Run-Tests
        }
        "fmt" {
            $Success = Format-Code
        }
        "vet" {
            $Success = Run-Vet
        }
        "deps" {
            $Success = Install-Dependencies
        }
        "install" {
            $Success = Install-Locally
        }
        "dev" {
            Run-Dev
            $Success = $true
        }
        "help" {
            Show-Help
            $Success = $true
        }
        default {
            Write-Host "Unknown target: $Target" -ForegroundColor Red
            Show-Help
            exit 1
        }
    }

    if (-not $Success) {
        exit 1
    }
}
catch {
    Write-Host "Error: $_" -ForegroundColor Red
    exit 1
}
