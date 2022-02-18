import pytest
import hoss
from hoss.dataset import PERM_NONE, PERM_READ, PERM_READ_WRITE


class TestDataset:
    def test_set_user_permission(self, fixture_test_dataset):
        s, ns, ds = fixture_test_dataset

        assert len(ds.permissions) == 2
        assert "admin" in ds.permissions.keys()
        assert "admin-hoss-default-group" in ds.permissions.keys()
        assert ds.permissions['admin'] == 'rw'

        ds.set_user_permission("privileged", PERM_READ)

        ds = ns.get_dataset(ds.dataset_name)

        assert len(ds.permissions) == 3
        assert ds.permissions['privileged-hoss-default-group'] == 'r'

        ds.set_user_permission("privileged", PERM_READ_WRITE)

        ds = ns.get_dataset(ds.dataset_name)

        assert len(ds.permissions) == 3
        assert ds.permissions['privileged-hoss-default-group'] == 'rw'

        ds.set_user_permission("privileged", PERM_NONE)

        ds = ns.get_dataset(ds.dataset_name)

        assert len(ds.permissions) == 2

    def test_set_cannot_remove_admin_group(self, fixture_test_dataset):
        s, ns, ds = fixture_test_dataset

        assert len(ds.permissions) == 2

        with pytest.raises(hoss.HossException):
            ds.set_group_permission("admin", None)

    def test_set_group_permission(self, fixture_test_dataset, fixture_test_auth):
        s, ns, ds = fixture_test_dataset
        s, group_name = fixture_test_auth

        assert len(ds.permissions) == 2

        s.auth.create_group(group_name, "my test group")
        s.auth.add_user_to_group(group_name, "user")

        ds.set_group_permission(group_name, PERM_READ)

        ds = ns.get_dataset(ds.dataset_name)

        assert len(ds.permissions) == 3
        assert ds.permissions['test-group'] == 'r'

        ds.set_group_permission(group_name, PERM_READ_WRITE)
        ds = ns.get_dataset(ds.dataset_name)

        assert len(ds.permissions) == 3
        assert ds.permissions['test-group'] == 'rw'

        ds.set_group_permission(group_name, PERM_NONE)
        ds = ns.get_dataset(ds.dataset_name)

        assert len(ds.permissions) == 2

    def test_set_public_group(self, fixture_test_dataset):
        s, ns, ds = fixture_test_dataset

        assert len(ds.permissions) == 2

        with pytest.raises(hoss.HossException):
            ds.set_group_permission("public", PERM_READ_WRITE)

        ds.set_group_permission("public", PERM_READ)

        ds = ns.get_dataset(ds.dataset_name)

        assert len(ds.permissions) == 3
        assert ds.permissions['public'] == 'r'
