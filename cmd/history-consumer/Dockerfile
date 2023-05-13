# Using moderc sqlite module:
# https://pkg.go.dev/modernc.org/sqlite#section-readme
# This removes the C sqlite3 libaries from the image along with
# requiring CGO_ENABLED
FROM gcr.io/distroless/static

COPY history-consumer /history-consumer
ENTRYPOINT [ "/history-consumer" ]
CMD ["consume", "-v"]
