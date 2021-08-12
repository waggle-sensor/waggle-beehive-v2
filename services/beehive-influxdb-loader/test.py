import unittest
from main import assert_type, assert_maxlen, assert_valid_message, assert_valid_message

class TestValidators(unittest.TestCase):

    def test_assert_type(self):
        assert_type(1, int)
        assert_type(1, (int, str, float))
        assert_type("wow", str)
        assert_type(1.23, float)
        with self.assertRaises(TypeError):
            assert_type("fail", int)
            assert_type(123, str)
    
    def test_assert_maxlen(self):
        assert_maxlen("123456789", 10)
        with self.assertRaises(ValueError):
            assert_maxlen("123456789", 8)

if __name__ == "__main__":
    unittest.main()
