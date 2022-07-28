function process(record) {
    var json = JSON.parse(String.fromCharCode.apply(String, record.Payload.Bytes()))
    if (json["trial"]) {
        return null;
    } else {
        return record
    }
}
