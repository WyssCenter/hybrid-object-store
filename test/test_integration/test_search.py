import hoss
import time
import pytest
from datetime import datetime as dt

import requests
import urllib.parse


class TestSearch:
    def test_search_basic(self, fixture_test_dataset):
        s, ns, ds = fixture_test_dataset

        f = ds / 'test1.txt'
        f.write_text("dummy data")

        f2 = ds / 'test2.txt'
        metadata = {"Foo": "bar", "fizz": "Buzz"}
        f2.write_text("more dummy data", metadata=metadata)

        f3 = ds / 'test3.txt'
        metadata = {"Fizz": "buzz"}
        f3.write_text("again dummy data", metadata=metadata)

        # Give time for indexing
        time.sleep(5)
        result = s.search({})
        assert len(result) == 3

        result = s.search({"Fizz": "buzz"})
        assert len(result) == 2

        result = s.search({"fizz": "Buzz"})
        assert len(result) == 2

        result = s.search({"fizz": "buzz"})
        assert len(result) == 2

        result = s.search({"Foo": "bar"})
        assert len(result) == 1

        result = s.search({"Foo": "buzz"})
        assert len(result) == 0

    def test_search_delete_object(self, fixture_test_dataset):
        s, ns, ds = fixture_test_dataset

        f1 = ds / 'test1.txt'
        metadata = {"Foo": "bar", "fizz": "Buzz"}
        f1.write_text("more dummy data", metadata=metadata)

        # Give time for indexing
        time.sleep(5)
        result = s.search({"fizz": "buzz"})
        assert len(result) == 1

        f1.remove()
        time.sleep(5)
        result = s.search({"fizz": "buzz"})
        assert len(result) == 0

    def test_search_recreate_object(self, fixture_test_dataset):
        s, ns, ds = fixture_test_dataset

        f1 = ds / 'test1.txt'
        metadata = {"foo": "bar", "meta": "witha:in:it"}
        f1.write_text("more dummy data", metadata=metadata)

        # Give time for indexing
        time.sleep(5)
        result = s.search({"meta": "witha:in:it"})
        assert len(result) == 1

        f1.remove()
        time.sleep(5)
        result = s.search({"meta": "witha:in:it"})
        assert len(result) == 0

        f1 = ds / 'test1.txt'
        metadata = {"new": "meta"}
        f1.write_text("recreate dummy data", metadata=metadata)
        time.sleep(5)

        f2 = ds / 'test1.txt'
        assert len(f2.metadata) == 1
        assert f2.metadata['new'] == 'meta'
        assert "foo" not in f2.metadata

        result = s.search({"meta": "witha:in:it"})
        assert len(result) == 0

        result = s.search({"new": "meta"})
        assert len(result) == 1

    def test_search_paging(self, fixture_test_dataset):
        s, ns, ds = fixture_test_dataset

        for i in range(26):
            file_ref = ds / f"path/to/file-{i}.ext"
            file_ref.write_text(f"test file {i}", metadata={"testing": 123})

        # Give time for indexing
        time.sleep(5)

        # test default page limit = 25, offset = 0
        result = s.search({"testing": 123})
        assert len(result) == 25

        # test setting higher offset
        result = s.search({"testing": 123}, offset=25)
        assert len(result) == 1

        # test setting smaller limit
        result = s.search({"testing": 123}, limit=10)
        assert len(result) == 10

        # test setting limit and offset
        result = s.search({"testing": 123}, limit=10, offset=20)
        assert len(result) == 6

        # test setting larger limit
        result = s.search({"testing": 123}, limit=30)
        assert len(result) == 26

        # test setting offset too high
        result = s.search({"testing": 123}, offset=30)
        assert len(result) == 0

    def test_search_refs(self, fixture_test_dataset):
        s, ns, ds = fixture_test_dataset

        f1 = ds / 'test1.txt'
        metadata = {"foo": "bar"}
        f1.write_text("dummy data", metadata=metadata)

        # Give time for indexing
        time.sleep(5)
        result = s.search_refs({"foo": "bar"})
        assert len(result) == 1
        assert result[0].dataset_name == ds.dataset_name

    def test_search_namespace_dataset_filters(self, fixture_test_datasets, fixture_test_alt_dataset):
        s, ns1, ds11, ds12 = fixture_test_datasets
        _, ns2, ds21 = fixture_test_alt_dataset

        f1 = ds11 / 'test1.txt'
        f2 = ds12 / 'test2.txt'
        f3 = ds21 / 'test3.txt'
        metadata = {"foo": "bar"}
        f1.write_text("dummy data", metadata=metadata)
        f2.write_text("more dummy data", metadata=metadata)
        f3.write_text("some more dummy data", metadata=metadata)

        # Give time for indexing
        time.sleep(5)

        # test that all come back without filtering
        result = s.search({"foo": "bar"})
        assert len(result) == 3

        # search by namespace
        result = s.search({"foo": "bar"}, namespace=ns1.name)
        assert len(result) == 2
        assert result[0]['namespace'] == ns1.name

        # search by dataset
        result = s.search({"foo": "bar"}, namespace=ns1.name, dataset=ds11.dataset_name)
        assert len(result) == 1
        assert result[0]['namespace'] == ns1.name
        assert result[0]['dataset'] == ds11.dataset_name

        # search by dataset without namespace (not allowed)
        with pytest.raises(hoss.HossException):
            result = s.search({"foo": "bar"}, dataset=ds11.dataset_name)

    def test_search_modified_filters(self, fixture_test_dataset):
        s, ns, ds = fixture_test_dataset

        metadata = {"foo": "bar"}
        timeFormat = "%Y-%m-%dT%H:%M:%S.%f"

        f1 = ds / f'test1.txt'
        f1.write_text(f"timed file 1", metadata=metadata)

        time.sleep(1)
        t1 = dt.utcnow().isoformat(timespec='milliseconds') + "Z"
        time.sleep(1)

        f2 = ds / f'test2.txt'
        f2.write_text(f"timed file 2", metadata=metadata)

        time.sleep(1)
        t2 = dt.utcnow().isoformat(timespec='milliseconds') + "Z"
        time.sleep(1)

        f3 = ds / f'test3.txt'
        f3.write_text(f"timed file 3", metadata=metadata)

        # Give time for indexing
        time.sleep(5)

        # test that all come back without filtering
        result = s.search({"foo": "bar"})
        assert len(result) == 3

        # test searching modified before
        result = s.search({"foo": "bar"}, modified_before=t2)
        assert len(result) == 2

        # test searching modified after
        result = s.search({"foo": "bar"}, modified_after=t1)
        assert len(result) == 2

        # test searching modified between
        result = s.search({"foo": "bar"}, modified_after=t1, modified_before=t2)
        assert len(result) == 1

        # test searching invalid range
        with pytest.raises(hoss.HossException):
            result = s.search({"foo": "bar"}, modified_after=t2, modified_before=t1)

    def test_search_key_autocomplete(self, fixture_test_datasets):
        s, ns1, ds11, ds12 = fixture_test_datasets

        f1 = ds11 / 'test1.txt'
        metadata = {"key1": "val1", "key2": "val1"}
        f1.write_text("dummy data", metadata=metadata)

        f2 = ds11 / 'test2.txt'
        metadata = {"key1": "val2", "key233": "val1"}
        f2.write_text("more dummy data", metadata=metadata)

        f3 = ds11 / 'test3.txt'
        metadata = {"key1": "val123", "key2345": "val1"}
        f3.write_text("some more dummy data", metadata=metadata)

        # Give time for indexing
        time.sleep(4)

        result = ds11.suggest_keys("key1")
        assert len(result) == 1
        assert result[0] == 'key1'

        result = ds11.suggest_keys("key2")
        assert len(result) == 3
        assert "key2" in result
        assert "key233" in result
        assert "key2345" in result

        result = ds11.suggest_keys("key2", limit=2)
        assert len(result) == 2

        result = ds11.suggest_keys("key23")
        assert len(result) == 2
        assert "key233" in result
        assert "key2345" in result

        result = ds11.suggest_keys("key234")
        assert len(result) == 1
        assert result[0] == "key2345"

    def test_search_value_autocomplete(self, fixture_test_datasets):
        s, ns1, ds11, ds12 = fixture_test_datasets

        f1 = ds11 / 'test1.txt'
        metadata = {"key1": "val1", "key2": "val1"}
        f1.write_text("dummy data", metadata=metadata)

        f2 = ds11 / 'test2.txt'
        metadata = {"key1": "val12", "key233": "val1"}
        f2.write_text("more dummy data", metadata=metadata)

        f3 = ds11 / 'test3.txt'
        metadata = {"key1": "val123", "key2345": "val1"}
        f3.write_text("some more dummy data", metadata=metadata)

        # Give time for indexing
        time.sleep(5)

        result = ds11.suggest_values("key1", "val")
        assert len(result) == 3
        assert "val1" in result
        assert "val12" in result
        assert "val123" in result

        result = ds11.suggest_values("key1", "val12")
        assert len(result) == 2
        assert "val12" in result
        assert "val123" in result

        result = ds11.suggest_values("key1", "val123")
        assert len(result) == 1
        assert "val123" in result

    def test_index_document_get(self, fixture_test_dataset):
        # There is an endpoint, primarily for the UI due to CORS, that can fetch the metadata of an object stored in the
        # search index. This test validates functionality, but since this isn't in the client library it's a bit manual
        s, ns, ds = fixture_test_dataset

        f1 = ds / 'simple-file.txt'
        metadata = {"test3": "val3"}
        f1.write_text("again dummy data", metadata=metadata)

        f2 = ds / "folder 1" / 'file with spaces.txt'
        metadata = {"test1": "val1", "test2": "val2"}
        f2.write_text("more dummy data", metadata=metadata)

        # Give time to index
        time.sleep(4)

        object_key = urllib.parse.quote_plus(f1.key)
        path = s.base_url + "/core/v1/search/namespace/" + ns.name + "/dataset/" + ds.dataset_name + "/metadata"
        payload = {"objectKey": object_key}

        response = requests.get(path, params=payload, headers={"Authorization": "Bearer " + s.auth.jwt})
        assert response.status_code == 200
        data = response.json()
        assert data['metadata']['test3'] == "val3"

        object_key = urllib.parse.quote_plus(f2.key)
        payload = {"objectKey": object_key}
        response = requests.get(path, params=payload, headers={"Authorization": "Bearer " + s.auth.jwt})
        assert response.status_code == 200
        data = response.json()
        assert data['metadata']['test1'] == "val1"
        assert data['metadata']['test2'] == "val2"

        object_key = urllib.parse.quote_plus(f2.key+"some junk")
        payload = {"objectKey": object_key}
        response = requests.get(path, params=payload, headers={"Authorization": "Bearer " + s.auth.jwt})
        assert response.status_code == 404
