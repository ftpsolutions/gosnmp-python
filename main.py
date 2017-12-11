from gosnmp_python import create_snmpv2c_session

session = create_snmpv2c_session(
    hostname='10.10.0.2',
    community='public',
)

session.connect()

print session.get('.1.3.6.1.2.1.1.5.0')

session.close()
