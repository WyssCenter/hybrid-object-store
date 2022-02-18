import pytest
import os
import shutil
import time

import hoss


def _clean_up_dataset(ns, dataset_name: str) -> None:
    """Helper function to delete a dataset and wait for it to be deleted in the background"""
    try:
        ns.delete_dataset(dataset_name)
    except hoss.error.NotFoundException:
        return

    while True:
        try:
            ns.get_dataset(dataset_name)
        except hoss.error.NotFoundException:
            return

        time.sleep(.5)


@pytest.fixture(scope="session")
def fixture_test_namespace():
    s = hoss.connect("http://localhost")

    src_bucket_path = os.path.expanduser("~/.hoss/data/nas/test-1")
    if os.path.isdir(src_bucket_path):
        shutil.rmtree(src_bucket_path, ignore_errors=True)
    os.makedirs(src_bucket_path)

    ns = s.create_namespace("test-src", "test source namespace", "default", "test-1")

    yield s, ns

    # clean up possible things that were created
    try:
        s.delete_namespace(ns.name)
    except hoss.error.NotFoundException:
        pass

    if os.path.isdir(src_bucket_path):
        shutil.rmtree(src_bucket_path, ignore_errors=True)

@pytest.fixture(scope="session")
def fixture_test_alt_namespace():
    s = hoss.connect("http://localhost")

    src_bucket_path = os.path.expanduser("~/.hoss/data/nas/test-2")
    if os.path.isdir(src_bucket_path):
        shutil.rmtree(src_bucket_path, ignore_errors=True)
    os.makedirs(src_bucket_path)

    ns = s.create_namespace("test-src-2", "another test source namespace", "default", "test-2")

    yield s, ns

    # clean up possible things that were created
    try:
        s.delete_namespace(ns.name)
    except hoss.error.NotFoundException:
        pass

    if os.path.isdir(src_bucket_path):
        shutil.rmtree(src_bucket_path, ignore_errors=True)


@pytest.fixture()
def fixture_test_namespace_name():
    yield "test_ns_2"

    # clean up possible things that were created
    try:
        s = hoss.connect("http://localhost")
        s.delete_namespace("test_ns_2")
    except hoss.error.NotFoundException:
        pass


@pytest.fixture()
def fixture_test_dataset_name(fixture_test_namespace):
    s, ns = fixture_test_namespace
    yield s, ns, "test_ds_1"

    # clean up possible things that were created
    _clean_up_dataset(ns, "test_ds_1")


@pytest.fixture()
def fixture_test_dataset(fixture_test_namespace):
    s, ns = fixture_test_namespace
    ds = ns.create_dataset("test_ds_1", "test dataset 1 in namespace 1")
    yield s, ns, ds

    # clean up possible things that were created
    _clean_up_dataset(ns, ds.dataset_name)


@pytest.fixture()
def fixture_test_alt_dataset(fixture_test_alt_namespace):
    s, ns = fixture_test_alt_namespace
    ds = ns.create_dataset("test_ds_3", "test dataset 1 in namespace 2")
    yield s, ns, ds

    # clean up possible things that were created
    _clean_up_dataset(ns, ds.dataset_name)


@pytest.fixture()
def fixture_test_datasets(fixture_test_namespace):
    s, ns = fixture_test_namespace
    ds1 = ns.create_dataset("test_ds_1", "test dataset 1 in namespace 1")
    ds2 = ns.create_dataset("test_ds_2", "test dataset 2 in namespace 1")
    yield s, ns, ds1, ds2

    # clean up possible things that were created
    _clean_up_dataset(ns, ds1.dataset_name)
    _clean_up_dataset(ns, ds2.dataset_name)


@pytest.fixture()
def fixture_test_dataset_with_data(fixture_test_dataset):
    s, ns, ds = fixture_test_dataset

    f = ds / "root.txt"
    f.write_text("dummy data - root.txt")

    f = ds / "folder1" / "foo1.txt"
    f.write_text("dummy data - folder1/foo1.txt")
    f = ds / "folder1" / "foo2.txt"
    f.write_text("dummy data - folder1/foo2.txt")
    f = ds / "folder1" / "bar1.txt"
    f.write_text("dummy data - folder1/bar1.txt")
    f = ds / "folder1" / "bar2.txt"
    f.write_text("dummy data - folder1/bar2.txt")

    yield s, ns, ds


@pytest.fixture()
def fixture_sync_config_local_local():
    s = hoss.connect("http://localhost")

    src_bucket_path = os.path.expanduser("~/.hoss/data/nas/source")
    if os.path.isdir(src_bucket_path):
        shutil.rmtree(src_bucket_path, ignore_errors=True)
    os.makedirs(src_bucket_path)

    ns_src = s.create_namespace("ns-src", "test source namespace", "default", "source")

    tgt_bucket_path = os.path.expanduser("~/.hoss/data/nas/target")
    if os.path.isdir(tgt_bucket_path):
        shutil.rmtree(tgt_bucket_path, ignore_errors=True)
    os.makedirs(tgt_bucket_path)

    ns_tgt = s.create_namespace("ns-tgt", "test target namespace", "default", "target")

    sync_config = ns_src.get_sync_configuration()
    assert sync_config['sync_enabled'] is False
    assert sync_config['sync_targets'] == list()

    ds_src_name = "ds-src"

    yield s, ns_src, ns_tgt, ds_src_name

    # Disable dataset sync if enabled
    try:
        ds = ns_src.get_dataset(ds_src_name)
        if ds.is_sync_enabled():
            ds.disable_sync()
            time.sleep(5)
    except hoss.error.NotFoundException:
        pass
    try:
        ds = ns_tgt.get_dataset(ds_src_name)
        if ds.is_sync_enabled():
            ds.disable_sync()
            time.sleep(5)
    except hoss.error.NotFoundException:
        pass

    ns_src.disable_sync_target("http://localhost", "ns-tgt")
    time.sleep(10)

    # clean up possible things that were created
    _clean_up_dataset(ns_src, ds_src_name)
    _clean_up_dataset(ns_tgt, ds_src_name)

    s.delete_namespace("ns-src")
    if os.path.isdir(src_bucket_path):
        shutil.rmtree(src_bucket_path, ignore_errors=True)
    s.delete_namespace("ns-tgt")
    if os.path.isdir(tgt_bucket_path):
        shutil.rmtree(tgt_bucket_path, ignore_errors=True)


@pytest.fixture()
def fixture_sync_config_local_remote():
    s = hoss.connect("http://localhost")

    local_bucket_path = os.path.expanduser("~/.hoss/data/nas/local")
    if os.path.isdir(local_bucket_path):
        shutil.rmtree(local_bucket_path, ignore_errors=True)
    os.makedirs(local_bucket_path)

    ns_local = s.create_namespace("ns-local", "test local namespace", "default", "local")
    ns_remote = s.create_namespace("ns-remote", "test remote namespace", "s3", "hos-int-test")

    sync_config = ns_local.get_sync_configuration()
    assert sync_config['sync_enabled'] is False
    assert sync_config['sync_targets'] == list()

    sync_config = ns_remote.get_sync_configuration()
    assert sync_config['sync_enabled'] is False
    assert sync_config['sync_targets'] == list()

    ds_src_name = "ds-src"

    yield s, ns_local, ns_remote, ds_src_name

    # Disable dataset sync if enabled
    try:
        ds = ns_local.get_dataset(ds_src_name)
        if ds.is_sync_enabled():
            ds.disable_sync()
            time.sleep(5)
    except hoss.error.NotFoundException:
        pass
    try:
        ds = ns_remote.get_dataset(ds_src_name)
        if ds.is_sync_enabled():
            ds.disable_sync()
            time.sleep(5)
    except hoss.error.NotFoundException:
        pass

    # Disable namespace sync if enabled
    sync_config = ns_local.get_sync_configuration()
    if sync_config['sync_enabled']:
        ns_local.disable_sync_target("http://localhost", "ns-remote")
        time.sleep(5)
    
    sync_config = ns_remote.get_sync_configuration()
    if sync_config['sync_enabled']:
        ns_remote.disable_sync_target("http://localhost", "ns-local")
        time.sleep(5)

    # clean up possible things that were created
    _clean_up_dataset(ns_local, ds_src_name)
    _clean_up_dataset(ns_remote, ds_src_name)

    s.delete_namespace("ns-local")
    if os.path.isdir(local_bucket_path):
        shutil.rmtree(local_bucket_path, ignore_errors=True)

    s.delete_namespace("ns-remote")
    time.sleep(2)


@pytest.fixture()
def fixture_test_auth():
    s = hoss.connect("http://localhost")
    test_group = "test-group"

    yield s, test_group

    try:
        s.auth.delete_group(test_group)
    except hoss.error.NotFoundException:
        pass


@pytest.fixture()
def fixture_test_public_group():
    s = hoss.connect("http://localhost")

    yield s

    s.auth.add_user_to_group("public", "user")



@pytest.fixture()
def fixture_default_policy() -> dict:
    return {"Version": "1", "Statements": []}


@pytest.fixture(scope="class")
def fixture_sync_policy_namespace():
    """Fixture for setting up namespaces when doing sync policy testing

    This fixture is used in the sync policy tests. It configures the namespace once for syncing. Since in these tests
    # we modify sync configs so much, this is one way of reducing test time a bit while not losing too much in test
    # case isolation. The `fixture_sync_policy_dataset` fixture is test scoped and does dataset clean up on each
    # test case
    """
    s = hoss.connect("http://localhost")

    src_bucket_path = os.path.expanduser("~/.hoss/data/nas/source")
    if os.path.isdir(src_bucket_path):
        shutil.rmtree(src_bucket_path, ignore_errors=True)
    os.makedirs(src_bucket_path)

    ns_src = s.create_namespace("ns-src", "test source namespace", "default", "source")

    tgt_bucket_path = os.path.expanduser("~/.hoss/data/nas/target")
    if os.path.isdir(tgt_bucket_path):
        shutil.rmtree(tgt_bucket_path, ignore_errors=True)
    os.makedirs(tgt_bucket_path)

    ns_tgt = s.create_namespace("ns-tgt", "test target namespace", "default", "target")

    sync_config = ns_src.get_sync_configuration()
    assert sync_config['sync_enabled'] is False
    assert sync_config['sync_targets'] == list()

    # Enable namespace syncing
    ns_src.enable_sync_target("http://localhost", "duplex", "ns-tgt")
    sync_config = ns_src.get_sync_configuration()
    assert sync_config['sync_enabled'] is True
    assert len(sync_config['sync_targets']) == 1
    assert sync_config['sync_targets'][0]['target_core_service'] == "http://localhost/core/v1"
    assert sync_config['sync_targets'][0]['target_namespace'] == "ns-tgt"
    assert sync_config['sync_targets'][0]['sync_type'] == "duplex"

    yield s, ns_src, ns_tgt

    ns_src.disable_sync_target("http://localhost", "ns-tgt")
    time.sleep(10)

    s.delete_namespace("ns-src")
    if os.path.isdir(src_bucket_path):
        shutil.rmtree(src_bucket_path, ignore_errors=True)
    s.delete_namespace("ns-tgt")
    if os.path.isdir(tgt_bucket_path):
        shutil.rmtree(tgt_bucket_path, ignore_errors=True)


@pytest.fixture()
def fixture_sync_policy_dataset(fixture_sync_policy_namespace):
    """Fixture for testing sync policies"""
    s, ns_src, ns_tgt = fixture_sync_policy_namespace

    ds_src_name = "ds-src-policy"

    yield s, ns_src, ns_tgt, ds_src_name

    # Disable dataset sync if enabled
    try:
        ds = ns_tgt.get_dataset(ds_src_name)
        if ds.is_sync_enabled():
            ds.disable_sync()
            time.sleep(8)
    except hoss.error.NotFoundException:
        pass
    try:
        ds = ns_src.get_dataset(ds_src_name)
        if ds.is_sync_enabled():
            ds.disable_sync()
            time.sleep(8)
    except hoss.error.NotFoundException:
        pass

    # clean up possible things that were created
    _clean_up_dataset(ns_src, ds_src_name)
    _clean_up_dataset(ns_tgt, ds_src_name)

    time.sleep(5)
