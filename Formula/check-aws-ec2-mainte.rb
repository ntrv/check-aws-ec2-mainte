class CheckAwsEc2Mainte < Formula
  desc 'Mackerel plugin to detect EC2 Retirement'
  homepage 'https://github.com/ntrv/check-aws-ec2-mainte'
  version '0.2.0'
  sha256 'undefined'
  url 'https://github.com/ntrv/check-aws-ec2-mainte/archive/v0.2.1.tar.gz'
  head 'https://github.com/ntrv/check-aws-ec2-mainte.git'

  def install
    bin.install 'check-aws-ec2-mainte'
  end

  test do
    system "#{bin}/check-aws-ec2-mainte", "--version"
  end
end
