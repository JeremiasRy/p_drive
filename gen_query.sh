if [ -f .backend.env ]; then
  source ./.backend.env
else
  echo ".backend.env file not found!"
  exit 1
fi

pass=$POSTGRES_PASSWORD
user=$POSTGRES_USER

url="postgres://$user:$pass@localhost:5432/$user?sslmode=disable"

jet -dsn=$url -path=./.gen