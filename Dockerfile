FROM openjdk:11

COPY target/ercole-server-*.jar /opt/ercole-server.jar

CMD java -Xms256M -Xmx256M -Dspring.profiles.active=sviluppo -jar /opt/ercole-server.jar
