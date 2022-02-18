import hoss
import hoss.error
import time
import pytest


class TestSyncLocalToLocal:
    def test_sync_simplex(self, fixture_sync_config_local_local):
        _, ns_src, ns_tgt, ds_src_name = fixture_sync_config_local_local

        ds_src = ns_src.create_dataset(ds_src_name, "source dataset")
        assert ds_src.is_sync_enabled() is False
        assert ds_src.sync_type is None
        assert ds_src.sync_policy is None

        # Dataset should not exist in target yet
        with pytest.raises(hoss.error.NotFoundException):
            ns_tgt.get_dataset(ds_src_name)

        # Enable namespace syncing
        ns_src.enable_sync_target("http://localhost", "simplex", "ns-tgt")
        sync_config = ns_src.get_sync_configuration()
        assert sync_config['sync_enabled'] is True
        assert len(sync_config['sync_targets']) == 1
        assert sync_config['sync_targets'][0]['target_core_service'] == "http://localhost/core/v1"
        assert sync_config['sync_targets'][0]['target_namespace'] == "ns-tgt"
        assert sync_config['sync_targets'][0]['sync_type'] == "simplex"

        # Enable dataset sync
        ds_src.enable_sync("simplex")
        assert ds_src.is_sync_enabled() is True
        time.sleep(5)
        ds_src = ns_src.get_dataset(ds_src_name)
        assert ds_src.sync_type == "simplex"
        assert ds_src.sync_policy is not None
        assert ds_src.sync_policy.get("Version") == '1'
        assert len(ds_src.sync_policy.get("Statements")) == 0

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
        f_src.write_text("a file written in the source dataset", metadata={"key1": "value1", "key2": "value2"})
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
        assert f_tgt.metadata["key1"] == "value1"
        assert f_tgt.metadata["key2"] == "value2"

        # Make sure URL parsing is working OK. include spaces and plus sign.
        # Write a file in the source and watch it end up in the target
        f_src = ds_src / 'test + 1 (3).txt'
        f_tgt = ds_tgt / 'test + 1 (3).txt'

        assert f_src.exists() is False
        assert f_tgt.exists() is False
        f_src.write_text("a file written in the source dataset again.")
        for cnt in range(30):
            try:
                assert f_src.exists() is True
                assert f_tgt.exists() is True
                print(f"Objects correct in {cnt} seconds")
                break
            except AssertionError:
                time.sleep(1)
                continue

        assert f_src.read_text() == f_tgt.read_text() == "a file written in the source dataset again."

        # The opposite should not work with simplex syncing
        f_src = ds_src / 'test2.txt'
        f_tgt = ds_tgt / 'test2.txt'

        assert f_src.exists() is False
        assert f_tgt.exists() is False
        f_tgt.write_text("a file written in the target dataset")
        time.sleep(5)
        assert f_src.exists() is False
        assert f_tgt.exists() is True

    def test_sync_duplex(self, fixture_sync_config_local_local):
        _, ns_src, ns_tgt, ds_src_name = fixture_sync_config_local_local

        ds_src = ns_src.create_dataset(ds_src_name, "source dataset")
        assert ds_src.is_sync_enabled() is False

        # Dataset should not exist in target yet
        with pytest.raises(hoss.error.NotFoundException):
            ns_tgt.get_dataset(ds_src_name)

        # Enable namespace syncing
        ns_src.enable_sync_target("http://localhost", "duplex", "ns-tgt")
        time.sleep(5)
        sync_config = ns_src.get_sync_configuration()
        assert sync_config['sync_enabled'] is True
        assert len(sync_config['sync_targets']) == 1
        assert sync_config['sync_targets'][0]['target_core_service'] == "http://localhost/core/v1"
        assert sync_config['sync_targets'][0]['target_namespace'] == "ns-tgt"
        assert sync_config['sync_targets'][0]['sync_type'] == "duplex"

        # Enable dataset sync
        ds_src.enable_sync("duplex")
        assert ds_src.is_sync_enabled() is True
        ds_src = ns_src.get_dataset(ds_src.dataset_name)

        time.sleep(70)

        print('waiting for api event processing and demuxer reload timeout...')
        # Dataset should now exist in the target namespace
        ds_tgt = None
        for cnt in range(15):
            try:
                ds_tgt = ns_tgt.get_dataset(ds_src_name)
                print(f"Target Dataset ready in {5*cnt} seconds")
                break
            except (hoss.error.NotFoundException, hoss.error.HossException):
                time.sleep(5)
                continue

        if not ds_tgt:
            raise Exception("Failed to load target dataset.")
        time.sleep(10)

        # Write a file in the source and watch it end up in the target
        f_src = ds_src / 'test.txt'
        f_tgt = ds_tgt / 'test.txt'

        assert f_src.exists() is False
        assert f_tgt.exists() is False
        f_src.write_text("a file written in the source dataset", metadata={"key1": "value1", "write": "source"})
        for cnt in range(60):
            try:
                assert f_src.exists() is True
                assert f_tgt.exists() is True
                print(f"Objects correct in {cnt} seconds")
                time.sleep(2)
                break
            except AssertionError:
                time.sleep(1)
                continue

        assert f_src.read_text() == f_tgt.read_text()
        assert f_src.metadata["key1"] == "value1"
        assert f_src.metadata["write"] == "source"
        assert f_tgt.metadata["key1"] == "value1"
        assert f_tgt.metadata["write"] == "source"

        # The opposite should work
        f_src = ds_src / 'test2.txt'
        f_tgt = ds_tgt / 'test2.txt'

        assert f_src.exists() is False
        assert f_tgt.exists() is False
        f_tgt.write_text("a file written in the target dataset", metadata={"key2": "value2", "write": "target"})
        for cnt in range(30):
            try:
                assert f_src.exists() is True
                assert f_tgt.exists() is True
                print(f"Objects correct in {cnt} seconds")
                time.sleep(3)
                break
            except AssertionError:
                time.sleep(1)
                continue

        assert f_src.read_text() == f_tgt.read_text()
        assert f_src.metadata["key2"] == "value2"
        assert f_src.metadata["write"] == "target"
        assert f_tgt.metadata["key2"] == "value2"
        assert f_tgt.metadata["write"] == "target"

    def test_sync_simplex_existing(self, fixture_sync_config_local_local):
        _, ns_src, ns_tgt, ds_src_name = fixture_sync_config_local_local

        ds_src = ns_src.create_dataset(ds_src_name, "source dataset")
        assert ds_src.is_sync_enabled() is False
        assert ds_src.sync_type is None
        assert ds_src.sync_policy is None

        # Dataset should not exist in target yet
        with pytest.raises(hoss.error.NotFoundException):
            ns_tgt.get_dataset(ds_src_name)

        # Enable namespace syncing
        ns_src.enable_sync_target("http://localhost", "simplex", "ns-tgt")
        sync_config = ns_src.get_sync_configuration()
        assert sync_config['sync_enabled'] is True
        assert len(sync_config['sync_targets']) == 1
        assert sync_config['sync_targets'][0]['target_core_service'] == "http://localhost/core/v1"
        assert sync_config['sync_targets'][0]['target_namespace'] == "ns-tgt"
        assert sync_config['sync_targets'][0]['sync_type'] == "simplex"

        # Write some data
        f1_src = ds_src / 'test1.txt'
        f2_src = ds_src / 'test2.txt'
        assert f1_src.exists() is False
        assert f2_src.exists() is False
        f1_src.write_text("a file written in the source dataset", metadata={"key1": "value1", "key2": "value2"})
        f2_src.write_text("a file written in the source dataset", metadata={"key1": "value1", "key2": "value2"})

        # Enable dataset sync
        ds_src.enable_sync("simplex")
        assert ds_src.is_sync_enabled() is True
        time.sleep(5)
        ds_src = ns_src.get_dataset(ds_src_name)
        assert ds_src.sync_type == "simplex"
        assert ds_src.sync_policy is not None
        assert ds_src.sync_policy.get("Version") == '1'
        assert len(ds_src.sync_policy.get("Statements")) == 0

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

        # Files should already exist in the target
        f1_tgt = ds_tgt / 'test1.txt'
        f2_tgt = ds_tgt / 'test2.txt'
        for cnt in range(30):
            try:
                assert f1_tgt.exists() is True
                assert f2_tgt.exists() is True
                print(f"Objects correct in {cnt} seconds")
                break
            except AssertionError:
                time.sleep(1)
                continue

        assert f1_src.read_text() == f1_tgt.read_text()
        assert f2_src.read_text() == f2_tgt.read_text()
        assert f1_tgt.metadata["key1"] == "value1"
        assert f2_tgt.metadata["key2"] == "value2"
