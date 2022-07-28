// Parses the record payload as JSON
function parseAsJSON(record) {
    return JSON.parse(String.fromCharCode.apply(String, record.Payload.Bytes()))
}

function transform(record) {
    logger.Info().Msg("entering transform");

    let json = parseAsJSON(record);
    json["greeting"] = "hello!";
    logger.Info().Msgf("json: %v", json);

    record.Payload = RawData();
    record.Payload.Raw = JSON.stringify(json);

    logger.Info().Msg("exiting transform");
    return record;
}
