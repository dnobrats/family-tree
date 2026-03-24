#!/bin/bash

# Script để reset sequences về giá trị tiếp theo sau ID lớn nhất
# Sử dụng: ./scripts/reset_sequences.sh

# Load environment variables nếu có
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Database connection info
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_NAME=${DB_NAME:-genealogy}

echo "🔄 Resetting sequences for database: $DB_NAME"
echo "================================================"

# Reset person sequence
echo "📝 Resetting person_id_seq..."
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c \
    "SELECT setval('person_id_seq', COALESCE((SELECT MAX(id) FROM person), 1), true);"

# Reset clan sequence
echo "📝 Resetting clan_id_seq..."
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c \
    "SELECT setval('clan_id_seq', COALESCE((SELECT MAX(id) FROM clan), 1), true);"

# Reset spouse sequence
echo "📝 Resetting spouse_id_seq..."
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c \
    "SELECT setval('spouse_id_seq', COALESCE((SELECT MAX(id) FROM spouse), 1), true);"

# Show results
echo ""
echo "✅ Sequences reset complete!"
echo "================================================"
echo "Current status:"
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c \
    "SELECT 
        'person' as table_name,
        COALESCE((SELECT MAX(id) FROM person), 0) as max_id,
        currval('person_id_seq') + 1 as next_id
    UNION ALL
    SELECT 
        'clan' as table_name,
        COALESCE((SELECT MAX(id) FROM clan), 0) as max_id,
        currval('clan_id_seq') + 1 as next_id
    UNION ALL
    SELECT 
        'spouse' as table_name,
        COALESCE((SELECT MAX(id) FROM spouse), 0) as max_id,
        currval('spouse_id_seq') + 1 as next_id;"

echo ""
echo "🎉 Done! Next IDs will continue from the numbers shown above."
