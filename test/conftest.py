import pytest

pytest_plugins = ['test_integration.fixtures']


def pytest_addoption(parser):
    parser.addoption(
        "--s3", action="store_true", default=False, help="run s3 tests"
    )


def pytest_configure(config):
    config.addinivalue_line("markers", "s3: mark test as requiring s3")


def pytest_collection_modifyitems(config, items):
    if config.getoption("--s3"):
        # --s3 given in cli: do not skip s3 tests
        return

    skip_s3 = pytest.mark.skip(reason="need --s3 option to run")
    for item in items:
        if "s3" in item.keywords:
            item.add_marker(skip_s3)
