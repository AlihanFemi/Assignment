# This is for MAC
export GOOSE_MIGRATION_DIR="migrations"

export GOOSE_DRIVER="postgres"

export GOOSE_DBSTRING="user=postgres dbname=assignmentdb password=admin sslmode=disable"


# This is for Windows
$env:GOOSE_MIGRATION_DIR="migrations"

$env:GOOSE_DRIVER="postgres"

$env:GOOSE_DBSTRING="user=postgres dbname=assignmentdb password=admin sslmode=disable"



# You will need to copy your sql code that makes your tables to your migration file, delete your tables
# goose -dir="migrations" postgres "user=postgres dbname=TestDB password=admin sslmode=disable" up | this is if you don't want to set up enviromental variables
# call | goose up | and test if your app still works and if it works that means your migration is done.
