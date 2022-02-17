# PowerAttackPublishingAPI

API to discover Homebrew D&D content for Patreon users of Power Attack Publishing.

## Local Setup

**Prerequisites**


**Start DB**

```bash
export POSTGRES_HOST=homebrew-db
export POSTGRES_DB=homebrew
export POSTGRES_USER=max
export POSTGRES_PASSWORD=<postgres-password>
make db-init

export DATA_PATH=data/spells.csv
export TABLE_NAME=spells
make db-write
```

**Start API**

```bash
make api
```

**Make requests**

- Navigate to localhost:8080/login in browser
- Login to Patreon


**Locally explore data**

```bash
make db-conn
```
