# Quiver System Requirements Validation (SRV) Overview

## Purpose
The SRV system ensures Quiver runs reliably by validating:

- Operating System compatibility
- CPU architecture
- Memory and disk resources
- Required dependencies (binaries/services)

It provides **early detection** of system issues and outputs results as a **terminal table** and **JSON report**.

---

## Architecture Diagram

```text
+-------------------+
|   configs/requirements.yaml  |
+-------------------+
           |
           v
+-------------------+
|   Manager (SRV)   |
|-------------------|
| - LoadRules()     |
| - ValidateAll()   |
| - Validators:    |
|   OS, Memory,     |
|   Disk, Dependencies |
+-------------------+
           |
           v
+-------------------+
|   Output          |
|-------------------|
| - Terminal Table  |
| - JSON Report     |
+-------------------+
