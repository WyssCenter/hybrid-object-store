import hoss
import hoss.error
import time
import pytest


@pytest.mark.s3
class TestSyncLocalToRemote:
    def test_sync_simplex(self, fixture_sync_config_local_remote):
        _, ns_src, ns_tgt, ds_src_name = fixture_sync_config_local_remote

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
        assert sync_config['sync_targets'][0]['target_namespace'] == "ns-remote"
        assert sync_config['sync_targets'][0]['sync_type'] == "simplex"

        # Enable dataset sync
        ds_src.enable_sync("simplex")
        assert ds_src.is_sync_enabled() is True
        time.sleep(20)

        print('waiting for api event processing and demuxer reload timeout...')
        ds_tgt = None
        for cnt in range(15):
            try:
                ds_tgt = ns_tgt.get_dataset(ds_src_name)
                print(f"Target Dataset ready in {5 * cnt} seconds")
                break
            except (hoss.error.NotFoundException, hoss.error.HossException):
                time.sleep(5)
                continue

        if not ds_tgt:
            raise Exception("Failed to load target dataset.")

        # Write a file in the source and watch it end up in the target
        f_src = ds_src / 'test.txt'
        f_tgt = ds_tgt / 'test.txt'

        assert f_src.exists() is False
        assert f_tgt.exists() is False

        f_src.write_text("a file written in the source dataset", metadata={"key1": "value1", "write": "local"})
        for cnt in range(30):
            try:
                assert f_src.exists() is True
                assert f_tgt.exists() is True
                print(f"Objects correct in {cnt} seconds")
                break
            except AssertionError:
                time.sleep(1)
                continue

        assert f_src.read_text() == f_tgt.read_text()
        assert f_src.metadata["key1"] == "value1"
        assert f_src.metadata["write"] == "local"
        assert f_tgt.metadata["key1"] == "value1"
        assert f_tgt.metadata["write"] == "local"

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
        _, ns_src, ns_tgt, ds_src_name = fixture_sync_config_local_remote

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
        assert sync_config['sync_targets'][0]['target_namespace'] == "ns-remote"
        assert sync_config['sync_targets'][0]['sync_type'] == "duplex"

        # Enable dataset sync
        ds_src.enable_sync("duplex")
        assert ds_src.is_sync_enabled() is True
        ds_src = ns_src.get_dataset(ds_src.dataset_name)
        time.sleep(70)

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
        time.sleep(10)

        # Write a file in the source and watch it end up in the target
        f_src = ds_src / 'test.txt'
        f_tgt = ds_tgt / 'test.txt'

        assert f_src.exists() is False
        assert f_tgt.exists() is False
        f_src.write_text("a file written in the source dataset", metadata={"key1": "value1", "write": "local"})
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
        assert f_src.metadata["key1"] == "value1"
        assert f_src.metadata["write"] == "local"
        assert f_tgt.metadata["key1"] == "value1"
        assert f_tgt.metadata["write"] == "local"

        # The opposite should work
        f_src = ds_src / 'test2.txt'
        f_tgt = ds_tgt / 'test2.txt'

        assert f_src.exists() is False
        assert f_tgt.exists() is False
        f_tgt.write_text("a file written in the target dataset", metadata={"key2": "value2", "write": "cloud"})
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
        assert f_src.metadata["key2"] == "value2"
        assert f_src.metadata["write"] == "cloud"
        assert f_tgt.metadata["key2"] == "value2"
        assert f_tgt.metadata["write"] == "cloud"
