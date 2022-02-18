import hoss
import tempfile
import time
import pytest

from hoss import utilities


class TestDatasetRef:
    def test_from_uri(self, fixture_test_dataset):
        s, ns, ds = fixture_test_dataset

        f = ds / "a_dir" / 'test.txt'
        f.write_text("dummy data")

        expected_uri = "hoss+http://localhost:test-src:test_ds_1/a_dir/test.txt"

        assert f.uri == expected_uri

        f_uri = hoss.resolve(expected_uri)
        data = f_uri.read_text()

        assert data == "dummy data"

    def test_file_exist(self, fixture_test_dataset_with_data):
        s, ns, ds = fixture_test_dataset_with_data

        f_root = ds / "root.txt"
        assert f_root.exists() is True
        assert f_root.is_file() is True
        assert f_root.is_dir() is False

        f_root2 = ds / "root-renamed.txt"
        assert f_root2.exists() is False

        f_root.move(f_root2)

        assert f_root2.exists() is True
        assert f_root.exists() is False

    def test_file_unlink(self, fixture_test_dataset_with_data):
        s, ns, ds = fixture_test_dataset_with_data

        f_root = ds / "root.txt"
        assert f_root.exists() is True
        assert f_root.is_file() is True
        assert f_root.is_dir() is False

        f_root.unlink()
        assert f_root.exists() is False

        f_root.unlink(missing_ok=True)

        with pytest.raises(FileNotFoundError):
            f_root.unlink()

    def test_file_remove(self, fixture_test_dataset_with_data):
        s, ns, ds = fixture_test_dataset_with_data

        f_root = ds / "root.txt"
        assert f_root.exists() is True
        assert f_root.is_file() is True
        assert f_root.is_dir() is False

        f_root.remove()
        assert f_root.exists() is False

        f_root.remove()
        assert f_root.exists() is False

    def test_glob(self, fixture_test_dataset_with_data):
        s, ns, ds = fixture_test_dataset_with_data

        objs = [f for f in ds.glob("*.txt")]
        assert len(objs) == 1

        objs = [f for f in ds.glob("**/*.txt")]
        assert len(objs) == 5

        objs = [f for f in ds.rglob("*.txt")]
        assert len(objs) == 5

        folder = ds / "folder1"
        objs = [f for f in folder.glob("*.txt")]
        assert len(objs) == 4

        objs = [f for f in folder.glob("foo*.txt")]
        assert len(objs) == 2

        objs = [f for f in folder.glob("bar*")]
        assert len(objs) == 2
        assert "bar" in objs[0].name

    def test_iterdir(self, fixture_test_dataset_with_data):
        s, ns, ds = fixture_test_dataset_with_data

        objs = [f for f in ds.iterdir()]
        assert len(objs) == 2

        for o in objs:
            if o.is_dir():
                objs2 = [f for f in o.iterdir()]
                assert len(objs2) == 4

    def test_protected_file(self, fixture_test_dataset_with_data):
        s, ns, ds = fixture_test_dataset_with_data

        with pytest.raises(hoss.HossException):
            ref = ds / ".dataset.yaml"

        with pytest.raises(hoss.HossException):
            ref = ds / "my-dir" / ".dataset.yaml"

        for f in ds.iterdir():
            assert f.name != ".dataset.yaml", ".dataset.yaml should not be returned to the user."

    def test_open(self, fixture_test_dataset):
        s, ns, ds = fixture_test_dataset

        f = ds / "test.txt"
        with f.open('wt') as fh:
            fh.write("dummy data 1")

        obj = ds / "test.txt"
        assert "dummy data 1" == obj.read_text()

        with f.open('rt') as fh:
            data = fh.read()

        assert data == "dummy data 1"

    def test_write_from_with_metadata(self, fixture_test_dataset):
        s, ns, ds = fixture_test_dataset

        with tempfile.NamedTemporaryFile(mode='wt') as tf:
            tf.write("some stuff")
            tf.flush()

            obj = ds / "my-file.dat"
            metadata = {"cat": "1", "dog": "2"}

            # Write a file with metadata
            obj.write_from(tf.name, metadata=metadata)
            time.sleep(2)

            obj2 = ds / "my-file.dat"
            assert obj2.metadata['cat'] == "1"
            assert obj2.metadata['dog'] == "2"
            etag1 = obj2.etag
            modified = obj2.last_modified
            assert obj2.size_bytes == 10
            assert obj2.read_text() == "some stuff"

            # change the file
            tf.write("some more stuff")
            tf.flush()

            # write a file again, removing metadata
            # Sleep because modified time has single second resolution
            time.sleep(2)
            metadata = {"cat": "1"}
            obj.write_from(tf.name, metadata=metadata)
            time.sleep(2)

            obj3 = ds / "my-file.dat"
            assert len(obj3.metadata) == 1
            assert obj3.metadata['cat'] == "1"
            assert etag1 != obj3.etag
            assert obj3.last_modified > modified
            assert obj3.size_bytes == 25
            assert obj3.read_text() == 'some stuffsome more stuff'

    def test_copy_to(self, fixture_test_dataset):
        s, ns, ds = fixture_test_dataset

        f1 = ds / "a_dir" / 'test.txt'
        f1.write_text("dummy data")
        assert f1.exists() is True

        f2 = ds / "b_dir" / 'test.txt'
        assert f2.exists() is False

        f1.copy_to(f2)
        assert f1.exists() is True
        assert f2.exists() is True

        assert f1.read_text() == f2.read_text()

    def test_hash_utilities_small(self, fixture_test_dataset):
        s, ns, ds = fixture_test_dataset

        # Check a small file under the multipart upload threshold
        with tempfile.NamedTemporaryFile(mode='wt') as tf:
            tf.write("this is a small file")
            tf.flush()

            obj = ds / "small-file.dat"

            # Check file hash function works as expected
            assert utilities.hash_file(tf.name) == "6955640f0bc822939f8bd440f32d62dd"

            # Write a file
            obj.write_from(tf.name)
            assert obj.etag == '"6955640f0bc822939f8bd440f32d62dd"'

            # Verify hash matches
            assert utilities.etag_does_match(obj, tf.name)

    def test_hash_utilities_large(self, fixture_test_dataset):
        s, ns, ds = fixture_test_dataset

        # Check a large file over the default multipart upload threshold (8MB)
        with tempfile.NamedTemporaryFile(mode='wt') as tf:
            tf.write("1234567890" * 1024 * 2000)
            tf.flush()

            obj = ds / "large-file.dat"

            # Check file hash function works as expected
            assert utilities.hash_file(tf.name) == '75948b7826fb8ee70bb057fb9997b873-3'

            # Write a file
            obj.write_from(tf.name)
            assert obj.etag == '"75948b7826fb8ee70bb057fb9997b873-3"'

            # Verify hash matches
            assert utilities.etag_does_match(obj, tf.name)
