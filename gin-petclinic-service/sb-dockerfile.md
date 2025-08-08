#syntax=cgregistry.capgroup.com/fix-container-certs
FROM cgregistry.capgroup.com/amazoncorretto:17-alpine-jdk AS build

USER root

ARG REVISION
ARG RELEASE_BRANCH

ARG ATM_ID
ARG PROJECT_PREFIX
ARG PROJECT_NAME

ARG NPM_CONFIG_REGISTRY

ARG DATADOG_JAVA_AGENT_ENABLED

ARG SONAR_SCAN_ENABLED
ARG SONAR_AUTH_TOKEN
ARG SONAR_BRANCH_NAME

ARG NEXUS_IQ_SCAN_ENABLED
ARG NEXUS_IQ_STAGE
ARG NEXUS_IQ_SERVER_URL
ARG NEXUS_IQ_API_ID
ARG NEXUS_IQ_API_KEY

ARG VERACODE_SCAN_ENABLED
ARG VERACODE_API_ID
ARG VERACODE_API_KEY

ARG CDP_OKTA_CLIENT_ID
ARG CDP_OKTA_CLIENT_SECRET
ARG CDP_OKTA_TOKEN_URI
ARG CDP_BASE_URL

ARG CHANNEL_SERVICE_CLIENT_ID
ARG CHANNEL_SERVICE_CLIENT_SECRET
ARG CHANNEL_SERVICE_TOKEN_URI
ARG CHANNEL_SERVICE_SCOPE
ARG CHANNEL_SERVICE_BASE_URL

ARG TEST_EMAIL_ADDRESS

ARG LAUNCHDARKLY_SDK_KEY

ARG DB_INTERACTION_SERVICE_BASE_URL
ARG OKTA_SYSTEM_ISSUER
ARG OKTA_INTERNAL_SERVICE_CLIENT_ID
ARG OKTA_INTERNAL_SERVICE_CLIENT_SECRET

ARG PLAID_CLIENT_ID
ARG PLAID_BASE_URL
ARG PLAID_CLIENT_SECRET

ARG TEST_EMAIL_ADDRESS

RUN apk add --no-cache python3 py3-pip npm git

RUN npm install -g @capgroup-cxt/cmn-build-tools@^0.1.0

COPY . /tmp/app-build
WORKDIR /tmp/app-build

RUN ./mvnw verify -V -P ci ${SONAR_SCAN_ENABLED:+"sonar:sonar"} -Drevision=${REVISION:-develop-SNAPSHOT}
-Dsonar.branch.name=$SONAR_BRANCH_NAME -Dsonar.login=$SONAR_AUTH_TOKEN
-Dveracode.api.id=$VERACODE_API_ID -Dveracode.api.key=$VERACODE_API_KEY

RUN if [ ! -z "$NEXUS_IQ_SCAN_ENABLED" ]; then \
nexus-iq scan target/*.jar --application-id ${ATM_ID}_$PROJECT_NAME --stage ${NEXUS_IQ_STAGE:-build}; \
fi

RUN if [ ! -z "$DATADOG_JAVA_AGENT_ENABLED" ]; then \
./mvnw dependency:copy -Dartifact=com.datadoghq:dd-java-agent:LATEST -DoutputDirectory=target/datadog
-Dmdep.stripClassifier=true -Dmdep.stripVersion=true; \
else \
mkdir -p target/datadog; \
fi

FROM cgregistry.capgroup.com/amazoncorretto:17-alpine

USER root

RUN mkdir -p /deployments/app /deployments/agents
RUN chown 0:root -R /deployments
RUN chmod -R 775 /deployments

WORKDIR /deployments/app

COPY --from=build /tmp/app-build/target/*.jar /deployments/app
COPY --from=build /tmp/app-build/target/datadog /deployments/agents

# Remove agents directory, if empty (i.e., Datadog agent does not exist)

RUN if [ ! "$(ls -A /deployments/agents)" ]; then \
rmdir /deployments/agents; \
fi

USER nobody
EXPOSE ${SERVER_PORT:-9080}

COPY docker-entrypoint.sh /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]