function transform(record) {
    logger.Info().Msg("entering transform")

    let jsonString = String.fromCharCode.apply(String, record.Payload.Bytes());
    logger.Info().Msgf("json string: %v", jsonString)

    var json = JSON.parse(jsonString)
    addGreeting(json)
    record.Payload.Raw = JSON.stringify(json);

    logger.Info().Msgf("payload: %+v", Object.keys(record.Payload))
    logger.Info().Msgf("jsonS: %v", json)
    return record;
}

function addGreeting(json) {
    json["greeting"] = "marhaba"
}