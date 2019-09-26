import os
import signal
import subprocess

from datetime import datetime

try:
    from subprocess import DEVNULL  # py3k
except ImportError:
    DEVNULL = open(os.devnull, 'wb')


class SNMPSimServer(object):
    """
    SNMPSim helper for starting a test snmp server. Useful for playing back a captured snmp walk.
    Capture snmp walks with the special snmp walk args (-ObentU) for compatibility with SNMP sim e.g.:

    snmpwalk -v2c -c public -ObentU localhost 1.3.6 > my_agent.snmpwalk
    """

    def __init__(self, data_dir, port):
        # You can run this command manually to start the snmpsimd server
        # Setting this to root - it's meant to be run in the docker container
        self._process = subprocess.Popen(
            'snmpsimd.py --process-user=root --process-group=root --data-dir={data_dir} --agent-udpv4-endpoint=127.0.0.1:{port}'.format(
                data_dir=data_dir,
                port=port),
            shell=True,
            env=os.environ,
            cwd=os.getcwd(),
            preexec_fn=os.setsid,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT)

        start_time = datetime.now()
        # Wait for the server to be ready...
        while True:
            output = ''
            line = self._process.stdout.readline()
            output += line
            if line:
                # The server is ready
                if line.strip().startswith('Listening at'):
                    break
            if (datetime.now() - start_time).total_seconds() > 10:
                raise Exception('Expected process to be ready by now but it isnt! output is: \n{}'.format(
                    output
                ))

    def close(self):
        os.killpg(os.getpgid(self._process.pid), signal.SIGTERM)

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        self.close()
        return False
