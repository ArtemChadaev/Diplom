---
description: It formalizes the current state by applying a version tag and performs a "Hard Reset" on the development environment.
---

Role: Version Manager

Step 1: Version Tagging
Identify Version: Check the last tag: git describe --tags --abbrev=0.

Ask the User: "Is this a **Release** (increment minor version, e.g., v0.1.0 -> v0.2.0) or a **Hotfix** (increment patch version, e.g., v0.2.0 -> v0.2.1)?"
**STOP and wait for the user's answer before proceeding to the next steps.**

Calculate Next Tag (v[Status].[Release].[Hotfix]) based on user's answer:

Release: Increment the middle digit and set the last digit to 0 (e.g., v0.1.0 -> v0.2.0).

Hotfix: Increment the last digit (e.g., v0.2.0 -> v0.2.1).

Apply: Create the tag: git tag -a vX.X.X -m "Finalized Release vX.X.X".

Step 2: Hard Reset of Development Branches
Reset Backend: >    * git checkout develop-backend

git reset --hard main

Reset Frontend:

git checkout develop-frontend

git reset --hard main

Warning: This will delete any uncommitted or divergent code in both development branches.

Step 3: Global Deployment
Push Main & Tags: git push origin main --tags.

Force Push Backend: git push origin develop-backend --force.

Force Push Frontend: git push origin develop-frontend --force.

Step 4: Completion Report
Confirm the new Version Tag.

Confirm that main, develop-backend, and develop-frontend are now perfectly synchronized across local and remote repositories.