export const getTopicsStatsMockData = () => {
  const data = [
    {
      name: 'publishers',
      value: 2,
    },
    {
      name: 'subscribers',
      value: 3,
    },
    {
      name: 'total_events',
      value: 1000000,
    },
    {
      name: 'storage',
      value: 203,
      units: 'MB',
    },
  ];
  return data;
};

export const getTopicEventsMockData = () => {
  return [
    {
      type: 'Document',
      version: '1.0.0',
      mimetype: 'application/json',
      events: {
        value: 12345678,
        percent: 96.0,
      },
      storage: {
        value: 512,
        units: 'MB',
        percent: 98.5,
      },
    },
    {
      type: 'Feed Item',
      version: '0.8.1',
      mimetype: 'application/rss',
      events: {
        value: 98765,
        percent: 4.0,
      },
      storage: {
        value: 4.3,
        units: 'KB',
        percent: 1.5,
      },
    },
  ];
};

export const createBinaryFixture = () => {
  const binaryData = new Uint8Array([72, 101, 108, 108, 111]);

  // Return the binary data as ArrayBuffer
  return binaryData.buffer;
};

export const getXMLFixture = () => {
  const xmlData = `
  <distributedSystem>
  <node id="node-1">
    <name>Node 1</name>
    <type>Application Server</type>
    <ipAddress>192.168.0.101</ipAddress>
    <port>8080</port>
    <status>Online</status>
    <connectedNodes>
      <nodeRef>node-2</nodeRef>
      <nodeRef>node-3</nodeRef>
    </connectedNodes>
  </node>

  <node id="node-2">
    <name>Node 2</name>
    <type>Database Server</type>
    <ipAddress>192.168.0.102</ipAddress>
    <port>3306</port>
    <status>Online</status>
    <connectedNodes>
      <nodeRef>node-1</nodeRef>
    </connectedNodes>
  </node>

  <node id="node-3">
    <name>Node 3</name>
    <type>Load Balancer</type>
    <ipAddress>192.168.0.103</ipAddress>
    <port>80</port>
    <status>Online</status>
    <connectedNodes>
      <nodeRef>node-1</nodeRef>
      <nodeRef>node-4</nodeRef>
      <nodeRef>node-5</nodeRef>
    </connectedNodes>
  </node>

  <node id="node-4">
    <name>Node 4</name>
    <type>Microservice</type>
    <ipAddress>192.168.0.104</ipAddress>
    <port>5000</port>
    <status>Offline</status>
    <connectedNodes>
      <nodeRef>node-3</nodeRef>
    </connectedNodes>
  </node>

  <node id="node-5">
    <name>Node 5</name>
    <type>Microservice</type>
    <ipAddress>192.168.0.105</ipAddress>
    <port>5000</port>
    <status>Online</status>
    <connectedNodes>
      <nodeRef>node-3</nodeRef>
    </connectedNodes>
  </node>
</distributedSystem>
`;
  return xmlData;
};

export const getTopicQueryResponseMockData = () => {
  return {
    results: [
      {
        metadata: {},
        mimetype: 'text/plain',
        version: 'Message v1.0.0',
        is_base64_encoded: false,
        data: 'hello world',
        created: 'seconds:1691428074  nanos:492983966',
      },
      {
        metadata: {},
        mimetype: 'text/csv',
        version: 'Spreadsheet v1.1.0',
        is_base64_encoded: false,
        data: 'hello,world',
        created: 'seconds:1691428074  nanos:493030132',
      },
      {
        metadata: {},
        mimetype: 'text/html',
        version: 'Webpage v1.0.1',
        is_base64_encoded: false,
        data: '<html><body><h1>Hello World</h1></body></html>',
        created: 'seconds:1691428074  nanos:493294049',
      },
      {
        metadata: {},
        mimetype: 'application/json',
        version: 'StockQuote v0.1.0',
        is_base64_encoded: false,
        data: '{\n  "price": 334.11,\n  "symbol": "MSFT",\n  "timestamp": 1690899527135,\n  "volume": 50\n}',
        created: 'seconds:1691428074  nanos:493561674',
      },
      {
        metadata: {},
        mimetype: 'application/json',
        version: 'StockQuote v0.1.0',
        is_base64_encoded: false,
        data: '{\n  "price": 320.23,\n  "symbol": "APPL",\n  "timestamp": 1690899527135,\n  "volume": 25\n}',
        created: 'seconds:1691428074  nanos:493562966',
      },
      {
        metadata: {},
        mimetype: 'application/json',
        version: 'StockQuote v0.1.0',
        is_base64_encoded: false,
        data: '{\n  "price": 335.12,\n  "symbol": "MSFT",\n  "timestamp": 1690899527135,\n  "volume": 40\n}',
        created: 'seconds:1691428074  nanos:493565382',
      },
      {
        metadata: {},
        mimetype: 'application/xml',
        version: 'Note v0.0.1',
        is_base64_encoded: true,
        data: 'PG5vdGU+PHRvPkFydGh1cjwvdG8+PGZyb20+TWFydmluPC9mcm9tPjxoZWFkaW5nPkxpZmU8L2hlYWRpbmc+PGJvZHk+RG9uJ3QgUGFuaWMhPC9ib2R5Pjwvbm90ZT4=',
        created: 'seconds:1691428074  nanos:493756716',
      },
      {
        metadata: {},
        mimetype: 'application/msgpack',
        version: 'Person v1.3.0',
        is_base64_encoded: true,
        data: 'gaRuYW1lo0JvYqNhZ2UYHg==',
        created: 'seconds:1691428074  nanos:493757841',
      },
      {
        metadata: {},
        mimetype: 'application/msgpack',
        version: 'Person v1.3.0',
        is_base64_encoded: true,
        data: 'gaRuYW1lo0FsaWNlo2FnZRge',
        created: 'seconds:1691428074  nanos:493758966',
      },
      {
        metadata: {},
        mimetype: 'application/protobuf',
        version: 'Person v2.0.1',
        is_base64_encoded: true,
        data: 'CgNCb2IQHg==',
        created: 'seconds:1691428074  nanos:493761174',
      },
    ],
    total_events: 11,
  };
};
