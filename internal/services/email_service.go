package services

import (
	"fmt"
	"log"

	"github.com/Habeebamoo/Clivo/server/internal/config"
	"github.com/Habeebamoo/Clivo/server/internal/models"
	"github.com/resend/resend-go/v3"
)

type EmailService interface {
	SendAdminMail(string, string) error
	SendWelcomeEmail(string, string, string)
	SendWelcomeEmailToAdmin(string, string, string, string)
	SendVerifiedUserEmail(string, string)
	SendUnverifiedUserEmail(string, string)
	SendRestrictedUserEmail(string, string)
	SendUnrestrictedUserEmail(string, string)
	SendCommentNotificationEmail(models.UserResponse, models.UserResponse, models.Article, models.Comment)
	SendCommentReplyNotificationEmail(models.UserResponse, models.UserResponse, models.Article, models.Comment)
}

type EmailSvc struct {}

func NewEmailService() EmailService {
	return &EmailSvc{}
}

func (ems *EmailSvc) SendAdminMail(name, email string) error {
	apiKey, _ := config.Get("RESEND_API_KEY")
	clientUrl, _ := config.Get("CLIENT_URL")

	imgUrl := fmt.Sprintf("%s/logo.png", clientUrl)

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
			<body style="font-family: Arial, sans-serif; background-color: #f4f4f7; padding: 10px; margin: 0;">
				<table width="100%%" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: auto; background: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 10px rgba(0,0,0,0.1);">
					<!-- logo -->
					<tr>
						<td style="padding: 20px; display: flex; justify-content: start; align-items: center; gap: 5px;">
							<img src="%s" style="height: 40px">
							<h1 style="font-family:'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif;">Clivo</h1>
						</td>
					</tr>

					<tr>
						<td style="padding: 0 30px; color: #333333; font-size: 16px; line-height: 1.6;">
							<p>Hi %s,</p>

							<p>
                You signed up for Clivo a few days ago - and i noticed you haven't published anything yet.
              </p>

              <p>
                That's totally normal. The blank page is the hardest part.
              </p>

              <p>
                So here's a low-pressure challenge: write anything. It dosen't have to be perfect. A short opinion, a how-to, a story from your week - just hit publish.
              </p>

              <p>
                Your first article is the hardest one. After that, it gets easier.
              </p>

              <p>
                &#x1F449; Start writing on Clivo.
              </p>

              <p>
                If you're stuck or have any questions, just reply to this email - I read every message.
              </p>

              <p>Cheers.</p>

							<div style="line-height: 0.4; margin-top: 30px;">
								<p>Habeeb Amoo</p>
								<p>Creator, Clivo.</p>
							</div>

							<p style="margin: 50px 0;">
								<a href="%s" style="background-color: rgb(20,20,20); color: #ffffff; padding: 12px 25px; border-radius: 5px; text-decoration: none; font-weight: bold;">Create Article</a>
							</p>
						</td>
					</tr>

					<tr>
						<td style="background-color: #f1f1f1; padding: 15px; text-align: center; font-size: 14px; color: #888888; line-height: 0;">
							<p style="font-family: Cambria, Cochin, Georgia, Times, 'Times New Roman', serif;">from</p>

							<div style="display: flex; align-items: center; gap: 3px; justify-content: center;">
								<img src="%s" style="height: 15px">
								<p style="font-family:'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif; color: black; font-weight: bold;">Clivo</p>  
							</div>

						</td>
					</tr>
				</table>
			</body>
		</html>
	`, imgUrl, name, clientUrl, imgUrl)

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From: "Habeeb from Clivo <hello@myclivo.com>",
		To: []string{email},
		Subject: "You joined Clivo - What's stopping you?",
		Html: html,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return err
	}
	
	log.Println("Email Request Received")
	return nil
}

func (ems *EmailSvc) SendWelcomeEmail(userName, userEmail, userUsername string) {
	apiKey, _ := config.Get("RESEND_API_KEY")
	clientUrl, _ := config.Get("CLIENT_URL")

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
			<body style="font-family: Arial, sans-serif; background-color: #f4f4f7; padding: 10px; margin: 0;">
				<table width="100%%" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: auto; background: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 10px rgba(0,0,0,0.1);">
					<!-- logo -->
					<tr>
						<td style="padding: 20px;">
							<img src="https://res.cloudinary.com/djvuchlcr/image/upload/c_fill,h_150,w_150/v1/profile_pics/fukp4ijlrcz9ojzrmy25?_a=AQAV6nF" style="height: 40px">
							<h1 style="font-family:'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif;">Clivo</h1>
						</td>
					</tr>

					<tr>
						<td style="padding: 0 30px; color: #333333; font-size: 16px; line-height: 1.6;">
							<p>Hi %s,</p>
							<p>Welcome to Clivo. we're excited to have you join our growing community of thinkers, writers, and readers.</p>

							<div style="line-height: 1.5;">
								<p>Here's what you can do next.</p>
								<p>&#9989; <span style="font-weight: bold;">Create</span> your first post and share your thoughts with the world.</p>
								<p>&#9989; <span style="font-weight: bold;">Discover</span> inspiring content from others who share your interests.</p>
								<p>&#9989; <span style="font-weight: bold;">Engage</span> with the community - like, comment and connect.</p>
							</div>

							<div style="line-height: 0.4; margin-top: 30px;">
								<p>Ready to start writing?</p>
								<p>Click the button below to create your first aticle.</p>
							</div>
							<p style="margin: 50px 0;">
								<a href="%s" style="background-color: rgb(20,20,20); color: #ffffff; padding: 12px 25px; border-radius: 5px; text-decoration: none; font-weight: bold;">Create Article</a>
							</p>
						</td>
					</tr>

					<tr>
						<td style="background-color: #f1f1f1; padding: 15px; text-align: center; font-size: 14px; color: #888888; line-height: 0;">
							<p style="font-family: Cambria, Cochin, Georgia, Times, 'Times New Roman', serif;">from</p>

							<div style="display: flex; align-items: center; gap: 3px; justify-content: center;">
								<img src="https://res.cloudinary.com/djvuchlcr/image/upload/c_fill,h_150,w_150/v1/profile_pics/fukp4ijlrcz9ojzrmy25?_a=AQAV6nF" style="height: 15px">
								<p style="font-family:'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif; color: black; font-weight: bold;">Clivo</p>  
							</div>

							<p style="font-family: Cambria, Cochin, Georgia, Times, 'Times New Roman', serif">
								This message was sent to 
								<span style="text-decoration: underline;">%s</span>
							</p>
						</td>
					</tr>
				</table>
			</body>
		</html>
	`, userName, fmt.Sprintf("%s/home/create", clientUrl), userEmail)

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From: "Clivo <hello@myclivo.com>",
		To: []string{userEmail},
		Subject: "Welcome To Clivo",
		Html: html,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("...Welcome Email Sent...")
}

//welcome email to admin
func (ems *EmailSvc) SendWelcomeEmailToAdmin(userName, userEmail, userUsername, interests string) {
	apiKey, _ := config.Get("RESEND_API_KEY")
	email, _ := config.Get("ADMIN_EMAIL")
	clientUrl, _ := config.Get("CLIENT_URL")

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
			<body style="font-family: Arial, sans-serif; background-color: #f4f4f7; padding: 10px; margin: 0;">
				<table width="100%%" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: auto; background: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 10px rgba(0,0,0,0.1);">
					<!-- logo -->
					<tr>
						<td style="padding: 20px;">
							<img src="https://res.cloudinary.com/djvuchlcr/image/upload/c_fill,h_150,w_150/v1/profile_pics/fukp4ijlrcz9ojzrmy25?_a=AQAV6nF" style="height: 40px">
							<h1 style="font-family:'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif;">Clivo</h1>
						</td>
					</tr>

					<tr>
						<td style="padding: 0 30px; color: #333333; font-size: 16px; line-height: 1.6;">
							<p>Good day Habeeb</p>
							<p>A new user recently signed up on Clivo. Below are the user details</p>

							<div style="line-height: 1.5;">
								<p>&#9989; <span style="font-weight: bold;">Name: </span>%s</p>
								<p>&#9989; <span style="font-weight: bold;">Email: </span>%s</p>
								<p>&#9989; <span style="font-weight: bold;">Username: </span>%s</p>
							</div>

							<p style="margin-top: 20px;">To view more about this user, click on the button below to visit the user's profile</p>

							<p style="margin: 50px 0;">
								<a href="%s style="background-color: rgb(20,20,20); color: #ffffff; padding: 12px 25px; border-radius: 5px; text-decoration: none; font-weight: bold;">View Profile</a>
							</p>
						</td>
					</tr>

					<tr>
						<td style="background-color: #f1f1f1; padding: 15px; text-align: center; font-size: 14px; color: #888888; line-height: 0;">
							<p style="font-family: Cambria, Cochin, Georgia, Times, 'Times New Roman', serif;">from</p>

							<div style="display: flex; align-items: center; gap: 3px; justify-content: center;">
								<img src="https://res.cloudinary.com/djvuchlcr/image/upload/c_fill,h_150,w_150/v1/profile_pics/fukp4ijlrcz9ojzrmy25?_a=AQAV6nF" style="height: 15px">
								<p style="font-family:'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif; color: black; font-weight: bold;">Clivo</p>  
							</div>
							<p style="font-family: Cambria, Cochin, Georgia, Times, 'Times New Roman', serif">
								<span>This message was sent to</span>
								<span style="text-decoration: underline;">habeebamoo08@gmail.com</span>
							</p>
						</td>
					</tr>
				</table>
			</body>
		</html>
	`, userName, userEmail, userUsername, fmt.Sprintf("%s/%s", clientUrl, userUsername))
	
	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From: "Clivo <hello@myclivo.com>",
		To: []string{email},
		Subject: "New User on Clivo",
		Html: html,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("...Welcome Email To Admin Sent...")
}

func (ems *EmailSvc) SendVerifiedUserEmail(userName, userEmail string) {
	apiKey, _ := config.Get("RESEND_API_KEY")

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
			<body style="font-family: Arial, sans-serif; background-color: #f4f4f7; padding: 10px; margin: 0;">
				<table width="100%%" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: auto; background: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 10px rgba(0,0,0,0.1);">
					<!-- logo -->
					<tr>
						<td style="padding: 20px;">
							<img src="https://res.cloudinary.com/djvuchlcr/image/upload/c_fill,h_150,w_150/v1/profile_pics/fukp4ijlrcz9ojzrmy25?_a=AQAV6nF" style="height: 40px">
							<h1 style="font-family:'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif;">Clivo</h1>
						</td>
					</tr>

					<tr>
						<td style="padding: 0 30px; color: #333333; font-size: 16px; line-height: 1.6;">
							<p>Hi %s,</p>
							<p>Great news! Your account has been officially verified, and you’ll now see a blue check mark next to your name on your profile and published articles.</p>

							<p>
                This verification lets readers know that your identity has been confirmed and that your work comes from an authentic, trusted voice on our platform. It also helps your articles stand out and builds credibility with your audience.
              </p>

              <p>
                There’s nothing you need to do — your verification is already live. Just keep writing, publishing, and engaging with the community as you always do.
              </p>

              <p>
                Thanks for being a valued part of our writer community. We’re excited to see what you publish next.
              </p>

							<div style="line-height: 0.4; margin-top: 30px;">
								<p>Warm regards</p>
								<p>The Clivo Team.</p>
							</div>
						</td>
					</tr>

					<tr>
						<td style="background-color: #f1f1f1; padding: 15px; text-align: center; font-size: 14px; color: #888888; line-height: 0;">
							<p style="font-family: Cambria, Cochin, Georgia, Times, 'Times New Roman', serif;">from</p>

							<div style="display: flex; align-items: center; gap: 3px; justify-content: center;">
								<img src="https://res.cloudinary.com/djvuchlcr/image/upload/c_fill,h_150,w_150/v1/profile_pics/fukp4ijlrcz9ojzrmy25?_a=AQAV6nF" style="height: 15px">
								<p style="font-family:'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif; color: black; font-weight: bold;">Clivo</p>  
							</div>

							<p style="font-family: Cambria, Cochin, Georgia, Times, 'Times New Roman', serif">
								This message was sent to 
								<span style="text-decoration: underline;">%s</span>
							</p>
						</td>
					</tr>
				</table>
			</body>
		</html>
	`, userName, userEmail)

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From: "Habeeb from Clivo <hello@myclivo.com>",
		To: []string{userEmail},
		Subject: "Your account is now verified",
		Html: html,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("...User Verified Email Sent...")
}

func (ems *EmailSvc) SendUnverifiedUserEmail(userName, userEmail string) {
	apiKey, _ := config.Get("RESEND_API_KEY")

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
			<body style="font-family: Arial, sans-serif; background-color: #f4f4f7; padding: 10px; margin: 0;">
				<table width="100%%" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: auto; background: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 10px rgba(0,0,0,0.1);">
					<!-- logo -->
					<tr>
						<td style="padding: 20px;">
							<img src="https://res.cloudinary.com/djvuchlcr/image/upload/c_fill,h_150,w_150/v1/profile_pics/fukp4ijlrcz9ojzrmy25?_a=AQAV6nF" style="height: 40px">
							<h1 style="font-family:'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif;">Clivo</h1>
						</td>
					</tr>

					<tr>
						<td style="padding: 0 30px; color: #333333; font-size: 16px; line-height: 1.6;">
							<p>Hi %s,</p>

							<p>
								We’re writing to let you know that the verification badge (blue check mark) has been removed from your account.
							</p>

							<p>
                This change doesn’t affect your ability to write, publish, or engage on the platform. Your profile and articles remain active and visible to readers as usual. Verification statuses may be updated from time to time based on our review processes and current eligibility criteria.
              </p>

              <p>
                If you believe this change was made in error or if you’d like to learn more about verification requirements, you can review our guidelines or reach out to our support team for assistance.
              </p>

              <p>
                Thank you for being part of our writer community, and we appreciate your continued contributions.
              </p>

							<div style="line-height: 0.4; margin-top: 30px;">
								<p>Warm regards</p>
								<p>The Clivo Team.</p>
							</div>
						</td>
					</tr>

					<tr>
						<td style="background-color: #f1f1f1; padding: 15px; text-align: center; font-size: 14px; color: #888888; line-height: 0;">
							<p style="font-family: Cambria, Cochin, Georgia, Times, 'Times New Roman', serif;">from</p>

							<div style="display: flex; align-items: center; gap: 3px; justify-content: center;">
								<img src="https://res.cloudinary.com/djvuchlcr/image/upload/c_fill,h_150,w_150/v1/profile_pics/fukp4ijlrcz9ojzrmy25?_a=AQAV6nF" style="height: 15px">
								<p style="font-family:'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif; color: black; font-weight: bold;">Clivo</p>  
							</div>

							<p style="font-family: Cambria, Cochin, Georgia, Times, 'Times New Roman', serif">
								This message was sent to 
								<span style="text-decoration: underline;">%s</span>
							</p>
						</td>
					</tr>
				</table>
			</body>
		</html>
	`, userName, userEmail)

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From: "Habeeb from Clivo <hello@myclivo.com>",
		To: []string{userEmail},
		Subject: "Update regarding your verification status",
		Html: html,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("...User Un-Verified Email Sent...")
}

func (ems *EmailSvc) SendRestrictedUserEmail(userName, userEmail string) {
	apiKey, _ := config.Get("RESEND_API_KEY")

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
			<body style="font-family: Arial, sans-serif; background-color: #f4f4f7; padding: 10px; margin: 0;">
				<table width="100%%" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: auto; background: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 10px rgba(0,0,0,0.1);">
					<!-- logo -->
					<tr>
						<td style="padding: 20px;">
							<img src="https://res.cloudinary.com/djvuchlcr/image/upload/c_fill,h_150,w_150/v1/profile_pics/fukp4ijlrcz9ojzrmy25?_a=AQAV6nF" style="height: 40px">
							<h1 style="font-family:'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif;">Clivo</h1>
						</td>
					</tr>

					<tr>
						<td style="padding: 0 30px; color: #333333; font-size: 16px; line-height: 1.6;">
							<p>Dear %s,</p>

							<p>
								We’re writing to inform you that your account has been temporarily restricted due to content that violates our Community Guidelines.
							</p>

							<p>
								Specifically, we identified one or more published articles that contain material inconsistent with our policies regarding [e.g., harassment, misinformation, hate speech, spam, copyright infringement]. As a result, your publishing privileges have been temporarily suspended.
							</p>

							<p>
								During this restriction period:
							</p>

							<ul>
								<li>You won't be able to login</li>
								<li>You will not be able to publish new articles.</li>
								<li>Existing content may be under review or permanently deleted.</li>
							</ul>

							<p>
								If you believe this action was taken in error, you may submit an appeal within 7 days by contacting our team at [clivoinc@gmail.com].
							</p>

							<p>
								We value thoughtful and responsible contributions to our platform and hope you will continue to share your voice in accordance with our policies.
							</p>

							<div style="line-height: 0.4; margin-top: 30px;">
								<p>Warm regards</p>
								<p>The Clivo Team.</p>
							</div>
						</td>
					</tr>

					<tr>
						<td style="background-color: #f1f1f1; padding: 15px; text-align: center; font-size: 14px; color: #888888; line-height: 0;">
							<p style="font-family: Cambria, Cochin, Georgia, Times, 'Times New Roman', serif;">from</p>

							<div style="display: flex; align-items: center; gap: 3px; justify-content: center;">
								<img src="https://res.cloudinary.com/djvuchlcr/image/upload/c_fill,h_150,w_150/v1/profile_pics/fukp4ijlrcz9ojzrmy25?_a=AQAV6nF" style="height: 15px">
								<p style="font-family:'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif; color: black; font-weight: bold;">Clivo</p>  
							</div>

							<p style="font-family: Cambria, Cochin, Georgia, Times, 'Times New Roman', serif">
								This message was sent to 
								<span style="text-decoration: underline;">%s</span>
							</p>
						</td>
					</tr>
				</table>
			</body>
		</html>
	`, userName, userEmail)

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From: "Habeeb from Clivo <hello@myclivo.com>",
		To: []string{userEmail},
		Subject: "Notice of Temporary Account Restriction",
		Html: html,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("...User Restriction Email Sent...")
}

func (ems *EmailSvc) SendUnrestrictedUserEmail(userName, userEmail string) {
	apiKey, _ := config.Get("RESEND_API_KEY")

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
			<body style="font-family: Arial, sans-serif; background-color: #f4f4f7; padding: 10px; margin: 0;">
				<table width="100%%" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: auto; background: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 10px rgba(0,0,0,0.1);">
					<!-- logo -->
					<tr>
						<td style="padding: 20px;">
							<img src="https://res.cloudinary.com/djvuchlcr/image/upload/c_fill,h_150,w_150/v1/profile_pics/fukp4ijlrcz9ojzrmy25?_a=AQAV6nF" style="height: 40px">
							<h1 style="font-family:'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif;">Clivo</h1>
						</td>
					</tr>

					<tr>
						<td style="padding: 0 30px; color: #333333; font-size: 16px; line-height: 1.6;">
							<p>Dear %s,</p>

							<p>
								We’re pleased to inform you that the restriction previously placed on your account has been lifted.
							</p>

							<p>
								After reviewing your account, we have restored full access to your profile and publishing privileges. You may now continue creating and engaging on the platform as usual.
							</p>

							<p>
								We appreciate your patience during the review process and your commitment to maintaining the standards of our community. If you have any questions or need clarification regarding our guidelines, please feel free to reach out to our support team.
							</p>

							<p>
								Thank you for being a valued member of our community.
							</p>

							<div style="line-height: 0.4; margin-top: 30px;">
								<p>Warm regards</p>
								<p>The Clivo Team.</p>
							</div>
						</td>
					</tr>

					<tr>
						<td style="background-color: #f1f1f1; padding: 15px; text-align: center; font-size: 14px; color: #888888; line-height: 0;">
							<p style="font-family: Cambria, Cochin, Georgia, Times, 'Times New Roman', serif;">from</p>

							<div style="display: flex; align-items: center; gap: 3px; justify-content: center;">
								<img src="https://res.cloudinary.com/djvuchlcr/image/upload/c_fill,h_150,w_150/v1/profile_pics/fukp4ijlrcz9ojzrmy25?_a=AQAV6nF" style="height: 15px">
								<p style="font-family:'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif; color: black; font-weight: bold;">Clivo</p>  
							</div>

							<p style="font-family: Cambria, Cochin, Georgia, Times, 'Times New Roman', serif">
								This message was sent to 
								<span style="text-decoration: underline;">%s</span>
							</p>
						</td>
					</tr>
				</table>
			</body>
		</html>
	`, userName, userEmail)

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From: "Habeeb from Clivo <hello@myclivo.com>",
		To: []string{userEmail},
		Subject: "Your Account Restriction Has Been Lifted",
		Html: html,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("...User Un-Restriction Email Sent...")
}

func (ems *EmailSvc) SendCommentNotificationEmail(author, commenter models.UserResponse, article models.Article, comment models.Comment) {
	apiKey, _ := config.Get("RESEND_API_KEY")
	clientUrl, _ := config.Get("CLIENT_URL")

	logoUrl := fmt.Sprintf("%s/logo.png", clientUrl)
	postUrl := fmt.Sprintf("%s/%s", clientUrl, article.Slug)
	subject := fmt.Sprintf("%s commented on your post", commenter.Name)

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>New Comment on Your Post</title>
			</head>
			<body style="font-family: Arial, sans-serif; background-color: #f4f4f7; padding: 10px; margin: 0; -webkit-font-smoothing: antialiased;">
				<table width="100%%" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: auto; background: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 10px rgba(0,0,0,0.05);">
					
					<!-- Header / Logo -->
					<tr>
						<td style="padding: 24px 30px 10px 30px;">
							<table width="100%%" cellspacing="0" cellpadding="0">
								<tr>
									<td width="40" style="vertical-align: middle;">
										<img src="%s" alt="Clivo Logo" style="height: 35px; display: block;">
									</td>
									<td style="padding-left: 8px; vertical-align: middle;">
										<span style="font-family:'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif; font-size: 22px; font-weight: bold; color: #111111;">Clivo</span>
									</td>
								</tr>
							</table>
						</td>
					</tr>

					<!-- Main Body -->
					<tr>
						<td style="padding: 20px 30px 30px 30px; color: #333333; font-size: 16px; line-height: 1.6;">
							<p style="margin-top: 0;">Hi %s,</p>

							<p style="margin-bottom: 24px;">
								<strong>%s</strong> just left a comment on your post, <span style="color: #111111; font-weight: 600;">"%s"</span>:
							</p>

							<!-- Comment Box Preview -->
							<table width="100%%" cellspacing="0" cellpadding="0" style="background-color: #f9f9fb; border-left: 4px solid #5d6ebd; border-radius: 4px; margin-bottom: 28px;">
								<tr>
									<td style="padding: 16px 20px; font-style: italic; color: #4a4a4a; font-size: 15px; line-height: 1.5;">
										"%s"
									</td>
								</tr>
							</table>

							<!-- Call to Action Button -->
							<table width="100%%" cellspacing="0" cellpadding="0" style="margin-bottom: 30px;">
								<tr>
									<td align="left">
										<a href="%s" target="_blank" style="background-color: #5d6ebd; color: #ffffff; padding: 12px 24px; font-weight: bold; text-decoration: none; border-radius: 6px; font-size: 15px; display: inline-block;">
											Reply to Comment
										</a>
									</td>
								</tr>
							</table>

							<p style="font-size: 14px; color: #666666; margin-bottom: 0;">
								If the button doesn't work, copy and paste this link into your browser:<br>
								<a href="%s" style="color: #5d6ebd; word-break: break-all; font-size: 13px;">%s</a>
							</p>
						</td>
					</tr>

					<!-- Footer -->
					<tr>
						<td style="background-color: #f8f9fa; padding: 20px; text-align: center; font-size: 12px; color: #888888; border-top: 1px solid #eeeeee;">
							<p style="font-family: Cambria, Georgia, serif; font-style: italic; margin: 0 0 8px 0; font-size: 14px;">from</p>
							
							<table align="center" cellspacing="0" cellpadding="0">
								<tr>
									<td style="vertical-align: middle; padding-right: 4px;">
										<img src="%s" alt="Clivo" style="height: 14px; display: block;">
									</td>
									<td style="vertical-align: middle;">
										<span style="font-family:'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif; color: #111111; font-weight: bold; font-size: 14px; letter-spacing: 0.5px;">Clivo</span>
									</td>
								</tr>
							</table>
						</td>
					</tr>
				</table>
			</body>
		</html>
	`, logoUrl, author.Name, commenter.Name, article.Title, comment.Content, postUrl, postUrl, postUrl, logoUrl)

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From: "Clivo <hello@myclivo.com>",
		To: []string{author.Email},
		Subject: subject,
		Html: html,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		log.Println(err)
	}
	
	log.Println("Email Request Received")
}

func (ems *EmailSvc) SendCommentReplyNotificationEmail(commenter, replier models.UserResponse, article models.Article, reply models.Comment) {
	apiKey, _ := config.Get("RESEND_API_KEY")
	clientUrl, _ := config.Get("CLIENT_URL")

	logoUrl := fmt.Sprintf("%s/logo.png", clientUrl)
	postUrl := fmt.Sprintf("%s/%s", clientUrl, article.Slug)
	subject := fmt.Sprintf("%s replied to your comment", replier.Name)

	html := fmt.Sprintf(`
    <!DOCTYPE html>
    <html lang="en">
      <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>New Reply to Your Comment</title>
      </head>
      <body style="font-family: Arial, sans-serif; background-color: #f4f4f7; padding: 10px; margin: 0; -webkit-font-smoothing: antialiased;">
        <table width="100%%" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: auto; background: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 10px rgba(0,0,0,0.05);">
          
          <!-- Header / Logo -->
          <tr>
            <td style="padding: 24px 30px 10px 30px;">
              <table width="100%%" cellspacing="0" cellpadding="0">
                <tr>
                  <td width="40" style="vertical-align: middle;">
                    <img src="%s" alt="Clivo Logo" style="height: 35px; display: block;">
                  </td>
                  <td style="padding-left: 8px; vertical-align: middle;">
                    <span style="font-family:'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif; font-size: 22px; font-weight: bold; color: #111111;">Clivo</span>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- Main Body -->
          <tr>
            <td style="padding: 20px 30px 30px 30px; color: #333333; font-size: 16px; line-height: 1.6;">
              <p style="margin-top: 0;">Hi %s,</p>

              <p style="margin-bottom: 24px;">
                <strong>%s</strong> just replied to your comment on the post, <span style="color: #111111; font-weight: 600;">"%s"</span>:
              </p>

              <!-- Reply Box Preview -->
              <table width="100%%" cellspacing="0" cellpadding="0" style="background-color: #f9f9fb; border-left: 4px solid #5d6ebd; border-radius: 4px; margin-bottom: 28px;">
                <tr>
                  <td style="padding: 16px 20px; font-style: italic; color: #4a4a4a; font-size: 15px; line-height: 1.5;">
                    "%s"
                  </td>
                </tr>
              </table>

              <!-- Call to Action Button -->
              <table width="100%%" cellspacing="0" cellpadding="0" style="margin-bottom: 30px;">
                <tr>
                  <td align="left">
                    <a href="%s" target="_blank" style="background-color: #5d6ebd; color: #ffffff; padding: 12px 24px; font-weight: bold; text-decoration: none; border-radius: 6px; font-size: 15px; display: inline-block;">
                      View Reply
                    </a>
                  </td>
                </tr>
              </table>

              <p style="font-size: 14px; color: #666666; margin-bottom: 0;">
                If the button doesn't work, copy and paste this link into your browser:<br>
                <a href="%s" style="color: #5d6ebd; word-break: break-all; font-size: 13px;">%s</a>
              </p>
            </td>
          </tr>

          <!-- Footer -->
          <tr>
            <td style="background-color: #f8f9fa; padding: 20px; text-align: center; font-size: 12px; color: #888888; border-top: 1px solid #eeeeee;">
              <p style="font-family: Cambria, Georgia, serif; font-style: italic; margin: 0 0 8px 0; font-size: 14px;">from</p>
              
              <table align="center" cellspacing="0" cellpadding="0">
                <tr>
                  <td style="vertical-align: middle; padding-right: 4px;">
                    <img src="%s" alt="Clivo" style="height: 14px; display: block;">
                  </td>
                  <td style="vertical-align: middle;">
                    <span style="font-family:'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif; color: #111111; font-weight: bold; font-size: 14px; letter-spacing: 0.5px;">Clivo</span>
                  </td>
                </tr>
              </table>
            </td>
          </tr>
        </table>
      </body>
    </html>
  `, logoUrl, commenter.Name, replier.Name, article.Title, reply.Content, postUrl, postUrl, postUrl, logoUrl)

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "Clivo <hello@myclivo.com>",
		To:      []string{commenter.Email},
		Subject: subject,
		Html:    html,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		log.Println(err)
	}

	log.Println("Reply Email Request Received")
}

