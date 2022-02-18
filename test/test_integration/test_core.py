import pytest
import os

import hoss
from hoss.error import HossException


class TestCoreService:
    def test_no_creds(self):
        pat = os.environ.get('HOSS_PAT')
        if not pat:
            raise Exception("YOU MUST SET THE `HOSS_PAT` env var for integration tests to run!")

        del os.environ['HOSS_PAT']

        with pytest.raises(HossException):
            hoss.connect("http://localhost")

        os.environ['HOSS_PAT'] = pat

    def test_list_object_stores(self, fixture_test_namespace):
        server, namespace = fixture_test_namespace

        obj_stores = server.list_object_stores()

        # at a minimum you have the default set up. you may have s3 too, but not enforcing here
        assert len(obj_stores) >= 1

    def test_get_object_store(self, fixture_test_namespace):
        server, namespace = fixture_test_namespace

        store = server.get_object_store("default")

        assert store.base_url == "http://localhost"
        assert store.description == "Default object store"
        assert store.endpoint == "http://localhost"
        assert store.host == "localhost"
        assert store.name == "default"
        assert store.object_store_type == "minio"

    def test_list_namespaces(self, fixture_test_namespace):
        server, _ = fixture_test_namespace

        namespaces = server.list_namespaces()

        # at a minimum you have the default set up
        assert len(namespaces) >= 1

    def test_get_namespace(self, fixture_test_namespace):
        server, _ = fixture_test_namespace

        ns = server.get_namespace("default")

        assert ns.name == "default"
        assert ns.host == "localhost"
        assert ns.description == "Default namespace"
        assert ns.bucket == "data"
        assert ns.base_url == "http://localhost"
        assert ns.auth is not None

    def test_create_delete_namespace(self, fixture_test_namespace, fixture_test_namespace_name):
        server, _ = fixture_test_namespace

        ns2 = server.create_namespace(fixture_test_namespace_name, "another ns", "default", "data")
        assert ns2.name == fixture_test_namespace_name
        assert ns2.host == "localhost"
        assert ns2.description == "another ns"
        assert ns2.bucket == "data"
        assert ns2.base_url == "http://localhost"
        assert ns2.auth is not None

        ns_get = server.get_namespace(fixture_test_namespace_name)
        assert ns_get.name == fixture_test_namespace_name
        assert ns_get.host == "localhost"
        assert ns_get.description == "another ns"
        assert ns_get.bucket == "data"
        assert ns_get.base_url == "http://localhost"
        assert ns_get.auth is not None

        server.delete_namespace(fixture_test_namespace_name)

        with pytest.raises(hoss.NotFoundException):
            server.get_namespace(fixture_test_namespace_name)
