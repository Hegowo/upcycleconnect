FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

ARG VITE_API_BASE_URL=/api/admin/v1
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

RUN npm run build

FROM httpd:alpine

COPY --from=builder /app/dist /usr/local/apache2/htdocs/
COPY apache.conf /usr/local/apache2/conf/extra/spa.conf

RUN sed -i \
    -e 's/#LoadModule rewrite_module/LoadModule rewrite_module/' \
    -e 's/#LoadModule expires_module/LoadModule expires_module/' \
    -e 's/#LoadModule headers_module/LoadModule headers_module/' \
    -e 's|#Include conf/extra/httpd-vhosts.conf|Include conf/extra/spa.conf|' \
    /usr/local/apache2/conf/httpd.conf

EXPOSE 80

CMD ["httpd-foreground"]
