#!/bin/sh

#DB=../../store/athletes/UHJvZmVzc29yIFpvb20\=.db
DB=$1
FINAL=$2
# ----------------
OUTPUT="results.csv"
SQL="suffer_query.sql"

# FINAL="score.json"

sqlite3 ${DB} \
        -init ${SQL} \
        ".exit" > ${OUTPUT}

echo "[" > ${FINAL}
cat ${OUTPUT} | \
        sort | \
        awk -F"|" '{print "{\"date\":\"" $1 "\",\"sst\":" $2 ",\"ss\":" $3 ",\"dist\":" $4 "},"}' \
        >> ${FINAL}
echo '{"date":"", "sst": 0, "ss": 0, "dist": 0} ]' >> ${FINAL}

rm ${OUTPUT}
