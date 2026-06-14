FROM httpd:2.4-alpine

RUN sed -i \
    -e 's|#LoadModule proxy_module|LoadModule proxy_module|' \
    -e 's|#LoadModule proxy_http_module|LoadModule proxy_http_module|' \
    -e 's|#LoadModule proxy_wstunnel_module|LoadModule proxy_wstunnel_module|' \
    -e 's|#LoadModule headers_module|LoadModule headers_module|' \
    -e 's|#LoadModule rewrite_module|LoadModule rewrite_module|' \
    -e 's|#LoadModule ssl_module|LoadModule ssl_module|' \
    -e 's|#LoadModule socache_shmcb_module|LoadModule socache_shmcb_module|' \
    /usr/local/apache2/conf/httpd.conf \
    && echo "ServerName localhost" >> /usr/local/apache2/conf/httpd.conf \
    && echo "Include conf/extra/vhost.conf" >> /usr/local/apache2/conf/httpd.conf \
    && printf '<IfFile /usr/local/apache2/conf/certs/fullchain.pem>\nListen 443\nSSLSessionCache "shmcb:/usr/local/apache2/logs/ssl_scache(512000)"\nInclude conf/extra/vhost-ssl.conf\n</IfFile>\n' >> /usr/local/apache2/conf/httpd.conf

COPY vhost.conf /usr/local/apache2/conf/extra/vhost.conf
COPY vhost-ssl.conf /usr/local/apache2/conf/extra/vhost-ssl.conf
