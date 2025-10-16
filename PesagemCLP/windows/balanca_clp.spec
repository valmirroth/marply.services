# -*- mode: python ; coding: utf-8 -*-
"""PyInstaller spec for BalancaCLP — generates an onedir Windows bundle.

Build from the project root with:
    pyinstaller --clean --noconfirm windows/balanca_clp.spec

Output: dist/BalancaCLP/BalancaCLP.exe (plus DLLs, templates, etc.)
"""
from pathlib import Path

ROOT = Path.cwd()
SNAP7_DLL = ROOT / "windows" / "snap7" / "snap7.dll"
TEMPLATES_DIR = ROOT / "templates"

if not SNAP7_DLL.exists():
    raise SystemExit(
        f"snap7.dll not found at {SNAP7_DLL}\n"
        "Download snap7-full-1.4.2 from SourceForge, extract release/Windows/Win64/snap7.dll\n"
        "and place it at windows/snap7/snap7.dll"
    )
if not TEMPLATES_DIR.is_dir():
    raise SystemExit(f"templates/ folder not found at {TEMPLATES_DIR}")


block_cipher = None

a = Analysis(
    [str(ROOT / "main.py")],
    pathex=[str(ROOT)],
    binaries=[
        (str(SNAP7_DLL), "."),
    ],
    datas=[
        (str(TEMPLATES_DIR), "templates"),
    ],
    hiddenimports=[
        "waitress",
        "pyodbc",
        "snap7",
        "snap7.client",
        "snap7.util",
    ],
    hookspath=[],
    runtime_hooks=[],
    excludes=[
        "tkinter",
        "matplotlib",
        "numpy",
        "pytest",
    ],
    win_no_prefer_redirects=False,
    win_private_assemblies=False,
    cipher=block_cipher,
)

pyz = PYZ(a.pure, a.zipped_data, cipher=block_cipher)

exe = EXE(
    pyz,
    a.scripts,
    [],
    exclude_binaries=True,
    name="BalancaCLP",
    debug=False,
    bootloader_ignore_signals=False,
    strip=False,
    upx=False,
    console=True,  # keep stdout/stderr visible — NSSM captures them
    icon=None,
)

coll = COLLECT(
    exe,
    a.binaries,
    a.zipfiles,
    a.datas,
    strip=False,
    upx=False,
    name="BalancaCLP",
)
