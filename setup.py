import setuptools
import subprocess
from setuptools import Distribution
from setuptools.command.build_py import build_py


class BinaryDistribution(Distribution):
    def has_ext_modules(foo):
        return True


with open("README.md", "r") as fh:
    long_description = fh.read()


class my_build_py(build_py):
    def run(self):
        # honor the --dry-run flag
        if not self.dry_run:
            return_value = subprocess.call(["./build.sh"])
            if return_value != 0:
                raise ValueError("build.sh returned non zero exit code")
        build_py.run(self)


setuptools.setup(
    name="gosnmp-python",
    version="0.91",
    # The project's main homepage.
    url="https://github.com/ftpsolutions/gosnmp-python",
    # Author details
    author="Edward @ FTP Technologies",
    author_email="edward.beech@ftpsolutions.com.au",
    # Choose your license
    license="MIT",
    description="GoSNMP Python",
    long_description=long_description,
    long_description_content_type="text/markdown",
    packages=setuptools.find_packages(),
    cmdclass={
        "build_py": my_build_py,
    },
    package_data={
        "": ["*.so"],
    },
    include_package_data=True,
    # Force the egg to unzip
    zip_safe=False,
    install_requires=["cffi==1.11.5", "future==0.17.1"],
    # Ensures that distributable copies are platform-specific and not universal
    distclass=BinaryDistribution,
    # See https://pypi.python.org/pypi?%3Aaction=list_classifiers
    classifiers=[
        # How mature is this project? Common values are
        #   3 - Alpha
        #   4 - Beta
        #   5 - Production/Stable
        "Development Status :: 7 - Inactive",
        # Indicate who your project is intended for
        "Intended Audience :: Developers",
        # Specify the Python versions you support here. In particular, ensure
        # that you indicate whether you support Python 2, Python 3 or both.
        "Programming Language :: Python :: 2.7",
        # OS Support
        "Operating System :: POSIX",
        "Operating System :: Unix",
        "Operating System :: MacOS",
    ],
)
