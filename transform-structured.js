function process(record) {
    record.Metadata["foo-key"] = "foo-value";
    return record;
}
