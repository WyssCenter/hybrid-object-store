import datetime
import pytest
import time

import hoss
from hoss.dataset import DeleteStatus


class TestNamespace:
    def test_sync_configuration(self, fixture_test_namespace, fixture_test_namespace_name):
        server, namespace = fixture_test_namespace

        server.create_namespace(fixture_test_namespace_name, "my test namespace to sync", "default", "data")

        namespace.enable_sync_target("http://localhost", "simplex", fixture_test_namespace_name)

        sync_config = namespace.get_sync_configuration()
        assert sync_config['sync_enabled'] is True
        assert len(sync_config['sync_targets']) == 1
        assert sync_config['sync_targets'][0]['target_core_service'] == 'http://localhost/core/v1'
        assert sync_config['sync_targets'][0]['target_namespace'] == fixture_test_namespace_name
        assert sync_config['sync_targets'][0]['sync_type'] == "simplex"

        namespace.disable_sync_target("http://localhost", fixture_test_namespace_name)

        sync_config = namespace.get_sync_configuration()
        assert sync_config['sync_enabled'] is False
        assert len(sync_config['sync_targets']) == 0

    def test_dataset_create_delete(self, fixture_test_dataset_name):
        server, namespace, ds_name = fixture_test_dataset_name

        ds = namespace.create_dataset(ds_name, "test dataset")
        assert ds.dataset_name == ds_name
        assert ds.description == "test dataset"
        assert ds.uri == f'hoss+http://localhost:{namespace.name}:{ds_name}/'
        assert ds.parent is None
        assert ds.bucket == 'test-1'
        assert ds.base_url == "http://localhost"
        assert ds.created_on < datetime.datetime.utcnow()
        assert ds.delete_status == DeleteStatus.NOT_SCHEDULED

        # Try to delete namespace, but fail because there is a dataset now
        with pytest.raises(hoss.HossException):
            server.delete_namespace(namespace.name)

        ds2 = namespace.get_dataset(ds_name)
        assert ds2.dataset_name == ds_name
        assert ds2.description == "test dataset"
        assert ds2.uri == f'hoss+http://localhost:{namespace.name}:{ds_name}/'
        assert ds2.parent is None
        assert ds2.bucket == 'test-1'
        assert ds2.base_url == "http://localhost"
        assert ds2.created_on < datetime.datetime.utcnow()
        assert ds.delete_status == DeleteStatus.NOT_SCHEDULED

        namespace.delete_dataset(ds_name)
        ds3 = namespace.get_dataset(ds_name)
        assert ds3.dataset_name == ds_name
        assert ds3.description == "test dataset"
        assert ds3.uri == f'hoss+http://localhost:{namespace.name}:{ds_name}/'
        assert ds3.parent is None
        assert ds3.bucket == 'test-1'
        assert ds3.base_url == "http://localhost"
        assert ds3.delete_status == DeleteStatus.SCHEDULED
        assert ds3.delete_on < datetime.datetime.utcnow()
        assert ds3.delete_on > datetime.datetime.utcnow() - datetime.timedelta(seconds=4)
        assert ds3.created_on < datetime.datetime.utcnow()

        time.sleep(4)
        with pytest.raises(hoss.NotFoundException):
            namespace.get_dataset(ds_name)

    def test_dataset_create_delete_restore(self, fixture_test_dataset_name):
        server, namespace, ds_name = fixture_test_dataset_name

        ds = namespace.create_dataset(ds_name, "test dataset")
        assert ds.dataset_name == ds_name
        assert ds.description == "test dataset"
        assert ds.uri == f'hoss+http://localhost:{namespace.name}:{ds_name}/'
        assert ds.parent is None
        assert ds.bucket == 'test-1'
        assert ds.base_url == "http://localhost"
        assert ds.created_on < datetime.datetime.utcnow()
        assert ds.delete_status == DeleteStatus.NOT_SCHEDULED

        namespace.delete_dataset(ds_name)
        ds3 = namespace.get_dataset(ds_name)
        assert ds3.dataset_name == ds_name
        assert ds3.delete_status == DeleteStatus.SCHEDULED

        namespace.restore_dataset(ds_name)

        ds3 = namespace.get_dataset(ds_name)
        assert ds3.dataset_name == ds_name
        assert ds3.delete_status == DeleteStatus.NOT_SCHEDULED

        # After 4 seconds it's still there (if you didn't successfully restore it would be deleted by now)
        time.sleep(4)
        ds2 = namespace.get_dataset(ds_name)
        assert ds2.dataset_name == ds_name
        assert ds2.description == "test dataset"
        assert ds2.uri == f'hoss+http://localhost:{namespace.name}:{ds_name}/'
        assert ds2.parent is None
        assert ds2.bucket == 'test-1'
        assert ds2.base_url == "http://localhost"
        assert ds2.created_on < datetime.datetime.utcnow()
        assert ds.delete_status == DeleteStatus.NOT_SCHEDULED

    def test_multi_group_dataset_list(self, fixture_test_dataset_name, fixture_test_auth):
        # This test makes sure that a user who is attached to a dataset multiple times
        # via multiple groups only has the dataset show up once (meaning deduplication
        # of permissions worked as expected)
        _, group_name = fixture_test_auth
        server, namespace, ds_name = fixture_test_dataset_name

        g1 = server.auth.create_group(group_name, "my test group")

        ds = namespace.create_dataset(ds_name, "test dataset")
        assert ds.dataset_name == ds_name
        assert ds.description == "test dataset"

        ds.set_group_permission(group_name, 'r')

        ds_list = namespace.list_datasets()

        assert len(ds_list) == 1

    def test_admin_sees_all_projects(self, fixture_test_dataset_name):
        # This test sort of "fakes" another user creating a dataset by adding the user with rw
        # and then remove admin. If you can still list and get the auto-admin group is what is
        # granting perms
        server, namespace, ds_name = fixture_test_dataset_name

        ds = namespace.create_dataset(ds_name, "test dataset")
        assert ds.dataset_name == ds_name
        assert ds.description == "test dataset"

        ds.set_user_permission("privileged", "rw")
        ds.set_user_permission("admin", None)

        ds_list = namespace.list_datasets()
        assert len(ds_list) == 1

        ds = namespace.get_dataset(ds_name)
        assert ds.dataset_name == ds_name

