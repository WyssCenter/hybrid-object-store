import hoss
import hoss.error
import time
import pytest

from test_integration.fixtures import fixture_sync_config_local_remote


@pytest.mark.s3
class TestSyncRemoteToLocal:
    def test_sync_simplex(self, fixture_sync_config_local_remote):
        _, ns_tgt, ns_src, ds_src_name = fixture_sync_config_local_remote

        ds_src = ns_src.create_dataset(ds_src_name, "source dataset")
        assert ds_src.is_sync_enabled() is False

        # Dataset should not exist in target yet
        with pytest.raises(hoss.error.NotFoundException):
            ns_tgt.get_dataset(ds_src_name)

        # Enable namespace syncing
        ns_src.enable_sync_target("http://localhost", "simplex", ns_tgt.name)
        sync_config = ns_src.get_sync_configuration()
        assert sync_config['sync_enabled'] is True
        assert len(sync_config['sync_targets']) == 1
        assert sync_config['sync_targets'][0]['target_core_service'] == "http://localhost/core/v1"
        assert sync_config['sync_targets'][0]['target_namespace'] == "ns-local"
        assert sync_config['sync_targets'][0]['sync_type'] == "simplex"

        # Enable dataset sync
        ds_src.enable_sync("simplex")
        assert ds_src.is_sync_enabled() is True

        # It can take a long time for S3 bucket events to start pumping through.
        time.sleep(30)

        ds_tgt = None
        for cnt in range(40):
            try:
                ds_tgt = ns_tgt.get_dataset(ds_src_name)
                print(f"Target Dataset ready in {2 * cnt} seconds")
                break
            except (hoss.error.NotFoundException, hoss.error.HossException):
                time.sleep(2)
                continue

        if not ds_tgt:
            raise Exception("Failed to load target dataset.")

        # Write a file in the source and watch it end up in the target
        f_src = ds_src / 'test.txt'
        f_tgt = ds_tgt / 'test.txt'

        assert f_src.exists() is False
        assert f_tgt.exists() is False
        time.sleep(5)
        f_src.write_text("a file written in the source dataset", metadata={"key4": "value4", "write": "cloud"})
        time.sleep(5)

        for cnt in range(3):
            try:
                assert f_src.exists() is True
                assert f_tgt.exists() is True
                print(f"Objects correct in {cnt} seconds")
                break
            except AssertionError:
                time.sleep(5)
                continue

        time.sleep(5)
        assert f_src.read_text() == f_tgt.read_text()
        assert f_src.metadata["key4"] == "value4"
        assert f_src.metadata["write"] == "cloud"
        assert f_tgt.metadata["key4"] == "value4"
        assert f_tgt.metadata["write"] == "cloud"

        # The opposite should not work with simplex syncing
        f_src = ds_src / 'test2.txt'
        f_tgt = ds_tgt / 'test2.txt'

        assert f_src.exists() is False
        assert f_tgt.exists() is False
        f_tgt.write_text("a file written in the target dataset")
        time.sleep(15)
        assert f_src.exists() is False
        assert f_tgt.exists() is True

    def test_sync_duplex(self, fixture_sync_config_local_remote):
        _, ns_tgt, ns_src, ds_src_name = fixture_sync_config_local_remote

        ds_src = ns_src.create_dataset(ds_src_name, "source dataset")
        assert ds_src.is_sync_enabled() is False

        # Dataset should not exist in target yet
        with pytest.raises(hoss.error.NotFoundException):
            ns_tgt.get_dataset(ds_src_name)

        # Enable namespace syncing
        ns_src.enable_sync_target("http://localhost", "duplex", ns_tgt.name)
        sync_config = ns_src.get_sync_configuration()
        assert sync_config['sync_enabled'] is True
        assert len(sync_config['sync_targets']) == 1
        assert sync_config['sync_targets'][0]['target_core_service'] == "http://localhost/core/v1"
        assert sync_config['sync_targets'][0]['target_namespace'] == "ns-local"
        assert sync_config['sync_targets'][0]['sync_type'] == "duplex"

        # Enable dataset sync
        ds_src.enable_sync("duplex")
        assert ds_src.is_sync_enabled() is True
        time.sleep(30)

        ds_src = ns_src.get_dataset(ds_src.dataset_name)
        # Dataset should now exist in the target namespace
        ds_tgt = None
        for cnt in range(30):
            try:
                ds_tgt = ns_tgt.get_dataset(ds_src_name)
                print(f"Target Dataset ready in {2*cnt} seconds")
                break
            except (hoss.error.NotFoundException, hoss.error.HossException):
                time.sleep(2)
                continue

        if not ds_tgt:
            raise Exception("Failed to load target dataset.")
        time.sleep(5)

        # Write a file in the source and watch it end up in the target
        f_src = ds_src / 'test.txt'
        f_tgt = ds_tgt / 'test.txt'

        assert f_src.exists() is False
        assert f_tgt.exists() is False
        f_src.write_text("a file written in the source dataset", metadata={"key5": "value5", "write": "cloud"})
        time.sleep(5)

        for cnt in range(30):
            try:
                assert f_src.exists() is True
                assert f_tgt.exists() is True
                print(f"Objects correct in {2*cnt} seconds")
                break
            except AssertionError:
                time.sleep(2)
                continue

        assert f_src.read_text() == f_tgt.read_text()
        assert f_src.metadata["key5"] == "value5"
        assert f_src.metadata["write"] == "cloud"
        assert f_tgt.metadata["key5"] == "value5"
        assert f_tgt.metadata["write"] == "cloud"

        # The opposite should work
        f_src = ds_src / 'test2.txt'
        f_tgt = ds_tgt / 'test2.txt'

        assert f_src.exists() is False
        assert f_tgt.exists() is False
        f_tgt.write_text("a file written in the target dataset", metadata={"key6": "value6", "write": "local"})
        time.sleep(10)
        f_src = ds_src / 'test2.txt'
        f_tgt = ds_tgt / 'test2.txt'

        for cnt in range(30):
            try:
                assert f_src.exists() is True
                assert f_tgt.exists() is True
                print(f"Objects correct in {2*cnt} seconds")
                break
            except AssertionError:
                time.sleep(2)
                continue

        assert f_src.read_text() == f_tgt.read_text()
        assert f_src.metadata["key6"] == "value6"
        assert f_src.metadata["write"] == "local"
        assert f_tgt.metadata["key6"] == "value6"
        assert f_tgt.metadata["write"] == "local"
