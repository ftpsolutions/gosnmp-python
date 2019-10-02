import traceback
import unittest
import os
from hamcrest import assert_that, equal_to
from gosnmp_python import create_snmpv2c_session
from concurrent.futures import ThreadPoolExecutor
import time

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

    def test_gosnmp_python_multi_threaded_get(self):

            def do_work():
                session = create_snmpv2c_session(
                            hostname='127.0.0.1',
                            community='dummy',
                            port=SIMULATOR_PORT
                        )
                try:
                    session.connect()
                    # Simulate several commands in sequence...
                    result1 = session.get('.1.3.6.1.2.1.1.1.0')
                    result2 = session.get('.1.3.6.1.2.1.1.1.0')
                    result3 = session.get('.1.3.6.1.2.1.1.1.0')
                    result4 = session.get('.1.3.6.1.2.1.1.1.0')
                    return '{}_{}_{}_{}'.format(result1.value, result2.value, result3.value, result4.value)
                except:
                    print(traceback.format_exc())
                finally:
                    session.close()

            futures = []

            num_iterations = 1000

            with ThreadPoolExecutor(max_workers=100) as executor:
                for _ in range(0, num_iterations):
                    futures.append(executor.submit(do_work))

            assert_that(len(futures), equal_to(num_iterations))

            for f in futures:
                assert_that(f.result(timeout=5), equal_to('Device Name_Device Name_Device Name_Device Name'))
