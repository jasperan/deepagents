#!/usr/bin/env bash
set -euo pipefail

ORACLE_HOST="${DEEPAGENTS_ORACLE_HOST:-localhost}"
ORACLE_PORT="${DEEPAGENTS_ORACLE_PORT:-1521}"
ORACLE_SERVICE="${DEEPAGENTS_ORACLE_SERVICE:-FREEPDB1}"
ORACLE_SYS_PWD="${DEEPAGENTS_ORACLE_PASSWORD:-DeepAgents2024}"
ORACLE_USER="${DEEPAGENTS_ORACLE_USER:-deepagents}"
ORACLE_USER_PWD="${DEEPAGENTS_ORACLE_PASSWORD:-DeepAgents2024}"

DSN="${ORACLE_HOST}:${ORACLE_PORT}/${ORACLE_SERVICE}"

echo "Waiting for Oracle to be ready at ${DSN}..."
for i in $(seq 1 60); do
    if echo "SELECT 1 FROM DUAL;" | sqlplus -s "sys/${ORACLE_SYS_PWD}@${DSN} as sysdba" &>/dev/null; then
        echo "Oracle is ready."
        break
    fi
    if [ "$i" -eq 60 ]; then
        echo "Timed out waiting for Oracle."
        exit 1
    fi
    sleep 5
done

echo "Creating user ${ORACLE_USER}..."
sqlplus -s "sys/${ORACLE_SYS_PWD}@${DSN} as sysdba" <<EOF
DECLARE
    user_exists INTEGER;
BEGIN
    SELECT COUNT(*) INTO user_exists FROM all_users WHERE username = UPPER('${ORACLE_USER}');
    IF user_exists = 0 THEN
        EXECUTE IMMEDIATE 'CREATE USER ${ORACLE_USER} IDENTIFIED BY "${ORACLE_USER_PWD}"';
        EXECUTE IMMEDIATE 'GRANT CONNECT, RESOURCE TO ${ORACLE_USER}';
        EXECUTE IMMEDIATE 'GRANT UNLIMITED TABLESPACE TO ${ORACLE_USER}';
        EXECUTE IMMEDIATE 'GRANT CREATE SESSION TO ${ORACLE_USER}';
    END IF;
END;
/
EOF

echo "Initializing schema..."
python3 -c "
import oracledb
conn = oracledb.connect(user='${ORACLE_USER}', password='${ORACLE_USER_PWD}', dsn='${DSN}')
from deepagents_oracle.schema import init_schema
init_schema(conn, vector=False)
conn.close()
print('Schema initialized.')
"

echo "Setup complete."
