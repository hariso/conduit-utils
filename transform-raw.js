function transform(record) {
    logger.Info().Msg("entering transform")

    var json = JSON.parse(String.fromCharCode.apply(String, record.Payload.Bytes))
    json["greeting"] = "marhaba"
    record.Payload.Raw = JSON.stringify(json);

    logger.Info().Msgf("payload: %+v", Object.keys(record.Payload))
    logger.Info().Msgf("jsonS: %v", json)
    return record;
}