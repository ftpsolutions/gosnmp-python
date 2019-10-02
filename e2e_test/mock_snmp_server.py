import os
import signal
import subprocess
import getpass
import threading

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
        self._stop = threading.Event()

        # If the user is root, as it is in docker containers, we need to handle that special case or
        # snmpsim will throw an error
        current_user = getpass.getuser()
        handle_root_user = '--process-user=root --process-group=root' if current_user == 'root' else ''

        # You can run this command manually to start the snmpsimd server
        # Setting this to root - it's meant to be run in the docker container
        self._process = subprocess.Popen(
            'snmpsimd.py {handle_root_user} --data-dir={data_dir} --agent-udpv4-endpoint=127.0.0.1:{port}'.format(
                data_dir=data_dir,
                port=port,
                handle_root_user=handle_root_user),
            shell=True,
            env=os.environ,
            cwd=os.getcwd(),
            preexec_fn=os.setsid,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT)

        start_time = datetime.now()
        # Wait for the server to be ready...
        output = ''
        while True:
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

        # Need to drain the process stdout because if the buffer gets full the subprocess stops running
        # properly ?!??
        self._drain_stdout_thread = threading.Thread(target=self._drain_stdout)
        self._drain_stdout_thread.start()

    def _drain_stdout(self):
        # Boolean setters are atomic
        while not self._stop.is_set():
            print(self._process.stdout.readline())

    def close(self):
        self._stop.set()
        os.killpg(os.getpgid(self._process.pid), signal.SIGTERM)

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        self.close()
        return False
