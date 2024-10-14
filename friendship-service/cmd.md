docker run --name myneo4j -p 7687:7687 -p 7474:7474 \
-e NEO4J_AUTH=neo4j/your_neo4j_password_here \
-e NEO4J_dbms_connector_bolt_advertised__address=localhost:7687 \
-e NEO4J_dbms_connector_http_advertised__address=localhost:7474 \
-d neo4j:latest