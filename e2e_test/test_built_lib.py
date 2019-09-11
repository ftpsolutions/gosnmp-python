import unittest
import os
from hamcrest import assert_that, equal_to
from gosnmp_python import create_snmpv2c_session

THIS_SCRIPT_DIR = os.path.dirname(os.path.realpath(__file__))
SIMULATOR_PORT = 4545

from mock_snmp_server import SNMPSimServer

class GoSNMPPythonLibTest(unittest.TestCase):

    def setUp(self):
        self._server = SNMPSimServer(data_dir=THIS_SCRIPT_DIR,
                                     port=SIMULATOR_PORT)
        self.addCleanup(self._server.close)

    def test_gosnmp_python_simple_get(self):
        session = create_snmpv2c_session(
            hostname='127.0.0.1',
            community='dummy',
            port=SIMULATOR_PORT
        )
        try:
            session.connect()
            result = session.get('.1.3.6.1.2.1.1.1.0')
            assert_that(result.value, equal_to('Device Name'))
        finally:
            session.close()
