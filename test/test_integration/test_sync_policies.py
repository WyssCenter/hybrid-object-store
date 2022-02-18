import hoss
import hoss.error
import time
import pytest
import copy


class TestSyncPolicies:
    def test_sync_policy_modify(self, fixture_sync_policy_dataset):
        _, ns_src, ns_tgt, ds_src_name = fixture_sync_policy_dataset

        ds_src = ns_src.create_dataset(ds_src_name, "source dataset")
        assert ds_src.is_sync_enabled() is False
        assert ds_src.sync_type is None
        assert ds_src.sync_policy is None

        # Enable dataset sync
        ds_src.enable_sync("simplex")
        assert ds_src.is_sync_enabled() is True
        time.sleep(2)
        ds_src = ns_src.get_dataset(ds_src_name)
        assert ds_src.sync_type == "simplex"
        assert ds_src.sync_policy is not None
        assert ds_src.sync_policy.get("Version") == '1'
        assert len(ds_src.sync_policy.get("Statements")) == 0

        statement = dict()
        statement["Id"] = "IgnoreRaw"
        statement["Conditions"] = [
                {
                    "Left": "object:key",
                    "Right": "*.raw",
                    "Operator": "!="
                }
            ]
        new_policy = ds_src.sync_policy
        new_policy['Statements'].append(statement)

        ds_src.enable_sync("simplex", sync_policy=new_policy)
        assert ds_src.is_sync_enabled() is True
        time.sleep(2)
        ds_src = ns_src.get_dataset(ds_src_name)
        assert ds_src.sync_type == "simplex"
        assert ds_src.sync_policy is not None
        assert ds_src.sync_policy.get("Version") == '1'
        assert len(ds_src.sync_policy.get("Statements")) == 1
        assert ds_src.sync_policy["Statements"][0]["Conditions"][0]["Left"] == "object:key"
        assert ds_src.sync_policy["Statements"][0]["Conditions"][0]["Right"] == "*.raw"
        assert ds_src.sync_policy["Statements"][0]["Conditions"][0]["Operator"] == "!="

    def test_sync_invalid_policy(self, fixture_sync_policy_dataset, fixture_default_policy):
        # A lot of tests case
        _, ns_src, ns_tgt, ds_src_name = fixture_sync_policy_dataset

        ds_src = ns_src.create_dataset(ds_src_name, "source dataset")
        assert ds_src.is_sync_enabled() is False
        assert ds_src.sync_type is None
        assert ds_src.sync_policy is None

        policy = fixture_default_policy
        statement = dict()
        statement["Id"] = "OnlyPuts"
        statement["Conditions"] = [
                {
                    "Left": "event:not-a-real-thing",
                    "Right": "PUT",
                    "Operator": "=="
                }
            ]
        policy['Statements'].append(statement)

        with pytest.raises(hoss.error.HossException):
            ds_src.enable_sync("duplex", sync_policy=policy)

    def test_sync_condition_operation(self, fixture_sync_policy_dataset, fixture_default_policy):
        _, ns_src, ns_tgt, ds_src_name = fixture_sync_policy_dataset

        ds_src = ns_src.create_dataset(ds_src_name, "source dataset")
        assert ds_src.is_sync_enabled() is False
        assert ds_src.sync_type is None
        assert ds_src.sync_policy is None

        policy = fixture_default_policy
        statement = dict()
        statement["Id"] = "OnlyPuts"
        statement["Conditions"] = [
                {
                    "Left": "event:operation",
                    "Right": "PUT",
                    "Operator": "=="
                }
            ]
        policy['Statements'].append(statement)

        ds_src.enable_sync("simplex", sync_policy=policy)
        assert ds_src.is_sync_enabled() is True
        time.sleep(15)
        ds_src = ns_src.get_dataset(ds_src_name)
        ds_tgt = ns_tgt.get_dataset(ds_src_name)

        f1_src = ds_src / "test1.txt"
        f1_src.write_text("test data 1")
        f2_src = ds_src / "test2.txt"
        f2_src.write_text("test data 2")

        time.sleep(5)

        f1_tgt = ds_tgt / "test1.txt"
        f2_tgt = ds_tgt / "test2.txt"

        assert f1_tgt.read_text() == f1_src.read_text()
        assert f2_tgt.read_text() == f2_tgt.read_text()

        f2_src.remove()

        assert f1_src.exists() is True
        assert f2_src.exists() is False
        assert f1_tgt.exists() is True
        assert f2_tgt.exists() is True

        f1_tgt.remove()

        assert f1_src.exists() is True
        assert f2_src.exists() is False
        assert f1_tgt.exists() is False
        assert f2_tgt.exists() is True

    def test_sync_condition_key(self, fixture_sync_policy_dataset, fixture_default_policy):
        _, ns_src, ns_tgt, ds_src_name = fixture_sync_policy_dataset

        ds_src = ns_src.create_dataset(ds_src_name, "source dataset")
        assert ds_src.is_sync_enabled() is False
        assert ds_src.sync_type is None
        assert ds_src.sync_policy is None

        policy = fixture_default_policy
        statement = dict()
        statement["Id"] = "NoDatFiles"
        statement["Conditions"] = [
                {
                    "Left": "object:key",
                    "Right": "*.dat",
                    "Operator": "!="
                }
            ]
        policy['Statements'].append(statement)

        ds_src.enable_sync("duplex", sync_policy=policy)
        assert ds_src.is_sync_enabled() is True
        time.sleep(75)
        ds_src = ns_src.get_dataset(ds_src_name)
        ds_tgt = ns_tgt.get_dataset(ds_src_name)

        f1_src = ds_src / "test1.txt"
        f1_src.write_text("test data 1")
        f2_src = ds_src / "test2.dat"
        f2_src.write_text("test data 2")

        time.sleep(10)

        f1_tgt = ds_tgt / "test1.txt"
        f2_tgt = ds_tgt / "test2.dat"

        assert f1_src.exists() is True
        assert f2_src.exists() is True
        assert f1_tgt.exists() is True
        assert f2_tgt.exists() is False
        assert f1_tgt.read_text() == f1_src.read_text()

        f1_tgt = ds_tgt / "test1-tgt.txt"
        f1_tgt.write_text("test data 1 written from target")
        f2_tgt = ds_tgt / "test2-tgt.dat"
        f2_tgt.write_text("test data 2 written from target")

        time.sleep(5)

        f1_tgt_src = ds_src / "test1-tgt.txt"
        f2_tgt_src = ds_src / "test2-tgt.dat"
        assert f1_tgt.exists() is True
        assert f2_tgt.exists() is True
        assert f1_tgt_src.exists() is True
        assert f2_tgt_src.exists() is False

    def test_sync_condition_size(self, fixture_sync_policy_dataset, fixture_default_policy):
        _, ns_src, ns_tgt, ds_src_name = fixture_sync_policy_dataset

        ds_src = ns_src.create_dataset(ds_src_name, "source dataset")
        assert ds_src.is_sync_enabled() is False
        assert ds_src.sync_type is None
        assert ds_src.sync_policy is None

        policy = fixture_default_policy
        statement = dict()
        statement["Id"] = "SmallFilesOnly"
        statement["Conditions"] = [
                {
                    "Left": "object:size",
                    "Right": 1000,
                    "Operator": "<"
                }
            ]
        policy['Statements'].append(statement)

        ds_src.enable_sync("simplex", sync_policy=policy)
        assert ds_src.is_sync_enabled() is True
        time.sleep(10)
        ds_src = ns_src.get_dataset(ds_src_name)
        ds_tgt = ns_tgt.get_dataset(ds_src_name)

        f1_src = ds_src / "small.txt"
        f1_src.write_text("test data 1")
        f2_src = ds_src / "big.txt"
        f2_src.write_text("1023456789" * 2000)

        time.sleep(8)

        f1_tgt = ds_tgt / "small.txt"
        f2_tgt = ds_tgt / "big.txt"

        assert f1_tgt.read_text() == f1_src.read_text()
        assert f1_src.exists() is True
        assert f2_src.exists() is True
        assert f1_tgt.exists() is True
        assert f2_tgt.exists() is False

    def test_sync_condition_metadata(self, fixture_sync_policy_dataset, fixture_default_policy):
        _, ns_src, ns_tgt, ds_src_name = fixture_sync_policy_dataset

        ds_src = ns_src.create_dataset(ds_src_name, "source dataset")
        assert ds_src.is_sync_enabled() is False
        assert ds_src.sync_type is None
        assert ds_src.sync_policy is None

        # Key has
        policy = copy.deepcopy(fixture_default_policy)
        statement = dict()
        statement["Id"] = "MustHaveTag"
        statement["Conditions"] = [
                {
                    "Left": "object:metadata",
                    "Right": "sync-me",
                    "Operator": "has"
                }
            ]
        policy['Statements'].append(statement)

        ds_src.enable_sync("simplex", sync_policy=policy)
        assert ds_src.is_sync_enabled() is True
        time.sleep(10)
        ds_src = ns_src.get_dataset(ds_src_name)
        ds_tgt = ns_tgt.get_dataset(ds_src_name)

        f1_src = ds_src / "test1.txt"
        f1_src.write_text("test data 1")
        f2_src = ds_src / "test2.txt"
        f2_src.write_text("test data 2", metadata={"sync-me": "1"})

        time.sleep(8)

        f1_tgt = ds_tgt / "test1.txt"
        f2_tgt = ds_tgt / "test2.txt"

        assert f2_tgt.read_text() == f2_src.read_text()
        assert f1_src.exists() is True
        assert f2_src.exists() is True
        assert f1_tgt.exists() is False
        assert f2_tgt.exists() is True

        # Value Match
        policy = copy.deepcopy(fixture_default_policy)
        statement = dict()
        statement["Id"] = "MatchValue"
        statement["Conditions"] = [
                {
                    "Left": "object:metadata:my-key",
                    "Right": "my-val-1",
                    "Operator": "=="
                }
            ]
        policy['Statements'].append(statement)

        ds_src.enable_sync("simplex", sync_policy=policy)
        assert ds_src.is_sync_enabled() is True
        time.sleep(70)
        ds_src = ns_src.get_dataset(ds_src_name)
        ds_tgt = ns_tgt.get_dataset(ds_src_name)

        f1_src = ds_src / "test3.txt"
        f1_src.write_text("test data 1", metadata={"my-key": "my-val-1", "other-key": "test"})
        f2_src = ds_src / "test4.txt"
        f2_src.write_text("test data 2", metadata={"my-key": "my-val-2"})

        time.sleep(8)

        f1_tgt = ds_tgt / "test3.txt"
        f2_tgt = ds_tgt / "test4.txt"

        assert f1_src.exists() is True
        assert f2_src.exists() is True
        assert f1_tgt.exists() is True
        assert f2_tgt.exists() is False
        assert f1_tgt.read_text() == f1_src.read_text()

        # Value Glob
        policy = copy.deepcopy(fixture_default_policy)
        statement = dict()
        statement["Id"] = "AnythingButMatchGlob"
        statement["Conditions"] = [
                {
                    "Left": "object:metadata:my-key",
                    "Right": "my-val*",
                    "Operator": "!="
                }
            ]
        policy['Statements'].append(statement)

        ds_src.enable_sync("simplex", sync_policy=policy)
        assert ds_src.is_sync_enabled() is True
        time.sleep(70)
        ds_src = ns_src.get_dataset(ds_src_name)
        ds_tgt = ns_tgt.get_dataset(ds_src_name)

        f1_src = ds_src / "test5.txt"
        f1_src.write_text("test data 1", metadata={"my-key": "my-val-1", "other-key": "test"})
        f2_src = ds_src / "test6.txt"
        f2_src.write_text("test data 2", metadata={"my-key": "my-val-2"})
        f3_src = ds_src / "test7.txt"
        f3_src.write_text("test data 3", metadata={"my-key": "other-val"})

        time.sleep(8)

        f1_tgt = ds_tgt / "test5.txt"
        f2_tgt = ds_tgt / "test6.txt"
        f3_tgt = ds_tgt / "test7.txt"

        assert f1_src.exists() is True
        assert f2_src.exists() is True
        assert f3_src.exists() is True
        assert f1_tgt.exists() is False
        assert f2_tgt.exists() is False
        assert f3_tgt.exists() is True
        assert f3_tgt.read_text() == f3_src.read_text()

    def test_sync_policy_multiple_conditions_and(self, fixture_sync_policy_dataset, fixture_default_policy):
        _, ns_src, ns_tgt, ds_src_name = fixture_sync_policy_dataset

        ds_src = ns_src.create_dataset(ds_src_name, "source dataset")
        assert ds_src.is_sync_enabled() is False
        assert ds_src.sync_type is None
        assert ds_src.sync_policy is None

        policy = fixture_default_policy
        statement = dict()
        statement["Id"] = "JustRightFilesOnly"
        statement["Effect"] = "AND"
        statement["Conditions"] = [
                {
                    "Left": "object:size",
                    "Right": 1000,
                    "Operator": "<"
                },
                {
                    "Left": "object:size",
                    "Right": 500,
                    "Operator": ">"
                }
            ]
        policy['Statements'].append(statement)

        ds_src.enable_sync("simplex", sync_policy=policy)
        assert ds_src.is_sync_enabled() is True
        time.sleep(10)
        ds_src = ns_src.get_dataset(ds_src_name)
        ds_tgt = ns_tgt.get_dataset(ds_src_name)

        f1_src = ds_src / "small.txt"
        f1_src.write_text("test data 1")
        f2_src = ds_src / "big.txt"
        f2_src.write_text("1023456789" * 3000)
        f3_src = ds_src / "right.txt"
        f3_src.write_text("1023456789" * 60)

        time.sleep(8)

        f1_tgt = ds_tgt / "small.txt"
        f2_tgt = ds_tgt / "big.txt"
        f3_tgt = ds_tgt / "right.txt"

        assert f1_src.exists() is True
        assert f2_src.exists() is True
        assert f3_src.exists() is True
        assert f1_tgt.exists() is False
        assert f2_tgt.exists() is False
        assert f3_tgt.exists() is True
        assert f3_tgt.read_text() == f3_src.read_text()

    def test_sync_policy_multiple_conditions_and_2(self, fixture_sync_policy_dataset, fixture_default_policy):
        _, ns_src, ns_tgt, ds_src_name = fixture_sync_policy_dataset

        ds_src = ns_src.create_dataset(ds_src_name, "source dataset")
        assert ds_src.is_sync_enabled() is False
        assert ds_src.sync_type is None
        assert ds_src.sync_policy is None

        policy = fixture_default_policy
        statement = dict()
        statement["Id"] = "JustRightFilesOnly"
        statement["Effect"] = "AND"
        statement["Conditions"] = [
            {
                "Left": "object:metadata",
                "Right": "key1",
                "Operator": "has"
            },
            {
                "Left": "object:metadata",
                "Right": "key2",
                "Operator": "has"
            }
        ]
        policy['Statements'].append(statement)

        ds_src.enable_sync("simplex", sync_policy=policy)
        assert ds_src.is_sync_enabled() is True
        time.sleep(10)
        ds_src = ns_src.get_dataset(ds_src_name)
        ds_tgt = ns_tgt.get_dataset(ds_src_name)

        f1_src = ds_src / "test1.txt"
        f1_src.write_text("test data 1", metadata={"key1": "val1", "key2": "val2"})
        f2_src = ds_src / "test2.txt"
        f2_src.write_text("test data 2", metadata={"key1": "val1"})
        f3_src = ds_src / "test3.txt"
        f3_src.write_text("test data 3", metadata={"key2": "val2"})

        time.sleep(8)

        f1_tgt = ds_tgt / "test1.txt"
        f2_tgt = ds_tgt / "test2.txt"
        f3_tgt = ds_tgt / "test3.txt"

        assert f1_src.exists() is True
        assert f2_src.exists() is True
        assert f3_src.exists() is True
        assert f1_tgt.exists() is True
        assert f2_tgt.exists() is False
        assert f3_tgt.exists() is False
        assert f1_tgt.read_text() == f1_src.read_text()

    def test_sync_policy_multiple_statements(self, fixture_sync_policy_dataset, fixture_default_policy):
        _, ns_src, ns_tgt, ds_src_name = fixture_sync_policy_dataset

        ds_src = ns_src.create_dataset(ds_src_name, "source dataset")
        assert ds_src.is_sync_enabled() is False
        assert ds_src.sync_type is None
        assert ds_src.sync_policy is None

        policy = fixture_default_policy
        statement = dict()
        statement["Id"] = "HaveKey"
        statement["Conditions"] = [
            {
                "Left": "object:metadata",
                "Right": "key1",
                "Operator": "has"
            }
        ]
        policy['Statements'].append(statement)

        statement = dict()
        statement["Id"] = "IsSmall"
        statement["Conditions"] = [
            {
                "Left": "object:size",
                "Right": 500,
                "Operator": "<"
            }
        ]
        policy['Statements'].append(statement)

        ds_src.enable_sync("simplex", sync_policy=policy)
        assert ds_src.is_sync_enabled() is True
        time.sleep(10)
        ds_src = ns_src.get_dataset(ds_src_name)
        ds_tgt = ns_tgt.get_dataset(ds_src_name)

        f1_src = ds_src / "test1.txt"
        f1_src.write_text("test data 1")
        f2_src = ds_src / "test2.txt"
        f2_src.write_text("0123456789" * 1000, metadata={"key1": "val1"})
        f3_src = ds_src / "test3.txt"
        f3_src.write_text("test data 3", metadata={"key1": "val1"})
        f4_src = ds_src / "test4.txt"
        f4_src.write_text("0123456789" * 1000)

        time.sleep(8)

        f1_tgt = ds_tgt / "test1.txt"
        f2_tgt = ds_tgt / "test2.txt"
        f3_tgt = ds_tgt / "test3.txt"
        f4_tgt = ds_tgt / "test4.txt"

        assert f1_src.exists() is True
        assert f2_src.exists() is True
        assert f3_src.exists() is True
        assert f4_src.exists() is True
        assert f1_tgt.exists() is True
        assert f2_tgt.exists() is True
        assert f3_tgt.exists() is True
        assert f4_tgt.exists() is False

        # Now make statements AND
        policy['Effect'] = "AND"

        ds_src.enable_sync("simplex", sync_policy=policy)
        assert ds_src.is_sync_enabled() is True
        time.sleep(70)
        ds_src = ns_src.get_dataset(ds_src_name)
        ds_tgt = ns_tgt.get_dataset(ds_src_name)

        f1_src = ds_src / "test5.txt"
        f1_src.write_text("test data 1")
        f2_src = ds_src / "test6.txt"
        f2_src.write_text("0123456789" * 1000, metadata={"key1": "val1"})
        f3_src = ds_src / "test7.txt"
        f3_src.write_text("test data 3", metadata={"key1": "val1"})
        f4_src = ds_src / "test8.txt"
        f4_src.write_text("0123456789" * 1000)

        time.sleep(8)

        f1_tgt = ds_tgt / "test5.txt"
        f2_tgt = ds_tgt / "test6.txt"
        f3_tgt = ds_tgt / "test7.txt"
        f4_tgt = ds_tgt / "test8.txt"

        assert f1_src.exists() is True
        assert f2_src.exists() is True
        assert f3_src.exists() is True
        assert f4_src.exists() is True
        assert f1_tgt.exists() is False
        assert f2_tgt.exists() is False
        assert f3_tgt.exists() is True
        assert f4_tgt.exists() is False

    def test_sync_update_policy(self, fixture_sync_policy_dataset, fixture_default_policy):
        _, ns_src, ns_tgt, ds_src_name = fixture_sync_policy_dataset

        ds_src = ns_src.create_dataset(ds_src_name, "source dataset")
        assert ds_src.is_sync_enabled() is False
        assert ds_src.sync_type is None
        assert ds_src.sync_policy is None

        policy = copy.deepcopy(fixture_default_policy)
        statement = dict()
        statement["Id"] = "NoDatFiles"
        statement["Conditions"] = [
                {
                    "Left": "object:key",
                    "Right": "*.dat",
                    "Operator": "!="
                }
            ]
        policy['Statements'].append(statement)

        ds_src.enable_sync("duplex", sync_policy=policy)
        assert ds_src.is_sync_enabled() is True
        time.sleep(75)
        ds_src = ns_src.get_dataset(ds_src_name)
        ds_tgt = ns_tgt.get_dataset(ds_src_name)

        f1_src = ds_src / "test1.txt"
        f1_src.write_text("test data 1")
        f2_src = ds_src / "test2.dat"
        f2_src.write_text("test data 2")

        time.sleep(10)

        f1_tgt = ds_tgt / "test1.txt"
        f2_tgt = ds_tgt / "test2.dat"

        assert f1_src.exists() is True
        assert f2_src.exists() is True
        assert f1_tgt.exists() is True
        assert f2_tgt.exists() is False
        assert f1_tgt.read_text() == f1_src.read_text()

        f1_tgt = ds_tgt / "test1-tgt.txt"
        f1_tgt.write_text("test data 1 written from target")
        f2_tgt = ds_tgt / "test2-tgt.dat"
        f2_tgt.write_text("test data 2 written from target")

        time.sleep(5)

        f1_tgt_src = ds_src / "test1-tgt.txt"
        f2_tgt_src = ds_src / "test2-tgt.dat"
        assert f1_tgt.exists() is True
        assert f2_tgt.exists() is True
        assert f1_tgt_src.exists() is True
        assert f2_tgt_src.exists() is False

        # Now switch the policy and make sure it is propagated to both namespaces
        policy = copy.deepcopy(fixture_default_policy)
        statement = dict()
        statement["Id"] = "NoTxtFiles"
        statement["Conditions"] = [
                {
                    "Left": "object:key",
                    "Right": "*.txt",
                    "Operator": "!="
                }
            ]
        policy['Statements'].append(statement)

        ds_src.enable_sync("duplex", sync_policy=policy)
        time.sleep(70)
        ds_src = ns_src.get_dataset(ds_src_name)
        ds_tgt = ns_tgt.get_dataset(ds_src_name)

        f1_src = ds_src / "test3.txt"
        f1_src.write_text("test data 3")
        f2_src = ds_src / "test4.dat"
        f2_src.write_text("test data 4")

        time.sleep(10)

        f1_tgt = ds_tgt / "test3.txt"
        f2_tgt = ds_tgt / "test4.dat"

        assert f1_src.exists() is True
        assert f2_src.exists() is True
        assert f1_tgt.exists() is False
        assert f2_tgt.exists() is True
        assert f2_tgt.read_text() == f2_src.read_text()

        f1_tgt = ds_tgt / "test3-tgt.txt"
        f1_tgt.write_text("test data 3 written from target")
        f2_tgt = ds_tgt / "test4-tgt.dat"
        f2_tgt.write_text("test data 4 written from target")

        time.sleep(5)

        f1_tgt_src = ds_src / "test3-tgt.txt"
        f2_tgt_src = ds_src / "test4-tgt.dat"
        assert f1_tgt.exists() is True
        assert f2_tgt.exists() is True
        assert f1_tgt_src.exists() is False
        assert f2_tgt_src.exists() is True
