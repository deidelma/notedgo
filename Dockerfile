FROM golang:1.20.5 as development 

WORKDIR /app 

RUN go install github.com/cosmtrek/air@latest

COPY go.mod ./ 
COPY go.sum ./
RUN go mod download 

# COPY . .

# RUN go install github.com/cespare/reflex@latest  

EXPOSE 5823

# CMD reflex -g ".go" go run main.go --start-service
CMD ["air", "-c", ".air.toml"]