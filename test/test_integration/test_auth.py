import pytest

import hoss
from hoss.auth import Role


class TestAuthService:
    def test_group(self, fixture_test_auth):
        server, group_name = fixture_test_auth

        g1 = server.auth.create_group(group_name, "my test group")
        assert g1.name == group_name
        assert g1.description == "my test group"
        assert len(g1.members) == 1
        assert g1.members[0].username == "admin"
        assert g1.members[0].role == Role.ADMIN

        g2 = server.auth.get_group(group_name)
        assert g2.name == group_name
        assert g2.description == "my test group"
        assert len(g2.members) == 1
        assert g2.members[0].username == "admin"
        assert g2.members[0].role == Role.ADMIN

        server.auth.delete_group(group_name)

        with pytest.raises(hoss.error.NotFoundException):
            server.auth.get_group(group_name)

    def test_group_users(self, fixture_test_auth):
        server, group_name = fixture_test_auth
        server.auth.create_group(group_name, "my test group")

        server.auth.add_user_to_group(group_name, "user")

        g1 = server.auth.get_group(group_name)
        assert g1.name == group_name
        assert g1.description == "my test group"
        assert len(g1.members) == 2

        server.auth.remove_user_from_group(group_name, "user")

        g1 = server.auth.get_group(group_name)
        assert g1.name == group_name
        assert g1.description == "my test group"
        assert len(g1.members) == 1

    def test_get_user(self, fixture_test_auth):
        server, group_name = fixture_test_auth

        u = server.auth.get_user('privileged')

        assert u.username == "privileged"
        assert u.email == "privileged@example.org"
        assert u.role == Role.PRIVILEGED

        with pytest.raises(hoss.error.NotFoundException):
            server.auth.get_user('doesnotexist')

    def test_check_default_groups(self, fixture_test_auth):
        server, group_name = fixture_test_auth

        admin_group = server.auth.get_group("admin")
        assert admin_group.name == "admin"
        assert len(admin_group.members) == 1

        public_group = server.auth.get_group("public")
        assert public_group.name == "public"
        assert len(public_group.members) == 3, "Public group not correct, did you login with all 3 test accounts?"

        with pytest.raises(hoss.HossException):
            server.auth.delete_group("admin")

        with pytest.raises(hoss.HossException):
            server.auth.delete_group("public")

    def test_admin_can_remove_user_from_public(self, fixture_test_public_group):
        server = fixture_test_public_group

        g1 = server.auth.get_group("public")
        assert len(g1.members) == 3, "Public group not correct, did you login with all 3 test accounts?"

        server.auth.remove_user_from_group("public", "user")

        g1 = server.auth.get_group("public")
        assert len(g1.members) == 2
