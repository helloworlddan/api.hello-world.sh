FROM gcr.io/endpoints-release/endpoints-runtime-serverless:2.37.0
USER root

ENV ESPv2_ARGS "^++^--cors_allow_credentials"
ENV ESPv2_ARGS "${ESPv2_ARGS}++--cors_preset=basic"
ENV ENDPOINTS_SERVICE_PATH /etc/endpoints/service.json
COPY service.json /etc/endpoints/service.json
RUN chown -R envoy:envoy /etc/endpoints/service.json && chmod -R 755 /etc/endpoints/service.json
USER envoy
ENTRYPOINT ["/env_start_proxy.py"]

# more config: https://cloud.google.com/endpoints/docs/openapi/specify-esp-v2-startup-options
