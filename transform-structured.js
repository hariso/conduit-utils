function transform(record) {
    logger.Info().Msg("entering transform");

    var p = StructuredData();
    p["newField"] = "marhaba, world!";
    record.Payload = p;

    logger.Info().Msgf("payload: %+v", Object.keys(record.Payload));
    logger.Info().Msgf("payload: %+v", Object.entries(record.Payload));
    logger.Info().Msgf("payload: %+v", record.Payload);

    logger.Info().Msgf("payload: %+v", Object.keys(p));
    logger.Info().Msgf("payload: %+v", Object.entries(p));
    logger.Info().Msgf("payload: %+v", p);

    logger.Info().Msg("exiting transform");
    return record;
}