# Usar la imagen base oficial de Go para compilar el código fuente
FROM golang:1.22 as builder

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar el código fuente local en el contenedor
COPY . .

# Compilar la aplicación Go.
# Desactivar el crosscompiling, habilitar el modo de módulos y compilar el binario.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o server

# Usar una imagen Docker minimalista para el contenedor final
FROM alpine:latest

# Instalar ca-certificates para llamadas HTTPS
RUN apk --no-cache add ca-certificates

# Establecer el directorio de trabajo en el contenedor
WORKDIR /root/

# Copiar el binario compilado desde el contenedor de compilación
COPY --from=builder /app/server .

# Exponer el puerto en el que la aplicación escuchará
EXPOSE 8080

# Ejecutar el binario compilado
CMD ["./server"]