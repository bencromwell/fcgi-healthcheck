<?php
header("Content-Type: text/plain");

echo "Hello from PHP-FPM!\n";
echo "REQUEST_METHOD: " . ($_SERVER["REQUEST_METHOD"] ?? "") . "\n";
echo "REQUEST_URI: " . ($_SERVER["REQUEST_URI"] ?? "") . "\n";
echo "QUERY_STRING: " . ($_SERVER["QUERY_STRING"] ?? "") . "\n";
echo "GET: " . json_encode($_GET) . "\n";
echo "POST: " . json_encode($_POST) . "\n";
$raw = file_get_contents('php://input');
echo "RAW: " . $raw . "\n";
