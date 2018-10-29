import setuptools
from setuptools.command.install import install
import subprocess


with open("README.md", "r") as fh:
    long_description = fh.read()

class CustomInstallCommand(install):
    """Customized setuptools install command - prints a friendly greeting."""
    def run(self):
        subprocess.call(['./build.sh'])
        install.run(self)


setuptools.setup(
    name="gosnmp-python",
    version="0.0.1",

    # The project's main homepage.
    url='https://github.com/ftpsolutions/gosnmp-python',

    # Author details
    author='scott @ FTP Technologies',
    author_email='scott.mills@ftpsolutions.com.au',

    # Choose your license
    license='Commercial',
    description="GoSNMP Python",
    long_description=long_description,
    long_description_content_type="text/markdown",
    packages=setuptools.find_packages(),
    cmdclass={'install': CustomInstallCommand},  # numpy hack
    package_data={
        'gosnmp-python': ['_gosnmp_python.so', 'gosnmp_python.so'],
    },
    include_package_data=True,

    install_requires=[
        'pytest',
        'mock',
        'pyhamcrest',
        'cffi',
        'future'
    ],

    # See https://pypi.python.org/pypi?%3Aaction=list_classifiers
    classifiers=[
        # How mature is this project? Common values are
        #   3 - Alpha
        #   4 - Beta
        #   5 - Production/Stable
        'Development Status :: 3 - Alpha',

        # Indicate who your project is intended for
        'Intended Audience :: Developers',
        'Topic :: FTP Technologies, IMS python tools',

        # Pick your license as you wish (should match "license" above)
        'Commercial',

        # Specify the Python versions you support here. In particular, ensure
        # that you indicate whether you support Python 2, Python 3 or both.
        'Programming Language :: Python :: 2.7',

        # OS Support
        "Operating System :: POSIX",
        "Operating System :: Unix",
        "Operating System :: MacOS",
    ],
)
