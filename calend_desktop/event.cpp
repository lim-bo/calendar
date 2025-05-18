#include "event.h"
#include "ui_event.h"

Event::Event(QString uid, QWidget *parent)
    : QWidget(parent)
    , ui(new Ui::Event), cli(settings.value("host").toString(), settings.value("port").toString()), uid(uid)
{
    ui->setupUi(this);
}

Event::~Event()
{
    delete ui;

}


void Event::on_create_event_clicked()
{
    ui->create_event->setDisabled(true);
    EventData event;
    event.master = uid;
    event.name = ui->nameEdit->text();
    event.desc = ui->descEdit->text();
    event.type = ui->typeCombo->currentText();
    event.prior = static_cast<EventData::Priority>(ui->priorityCombo->currentIndex() + 1);

    event.start = QDateTime(ui->startDateEdit->date(), ui->startTimeEdit->time());
    event.end = QDateTime(ui->endDateEdit->date(), ui->endTimeEdit->time());

    QStringList partsEmails = ui->parts->toPlainText().split(" ", Qt::SkipEmptyParts);
    for (const QString &email : partsEmails) {
        Participant participant;
        participant.uid = email;
        participant.accepted = false;
        event.parts.append(participant);
    }

    int notifyIndex = ui->notificationCombo->currentIndex();
    if (notifyIndex > 0) {
        event.notificationTime = calculateNotificationTime(event.start, notifyIndex);


    }

    QString eventID;
    bool ok;
    if (isEditMode) {
        // ok = cli.deleteEvent(uid, originalEventId);
        // if (ok) {
        //     ok = cli.addEvent(event, uid, &eventID);
        // }
        event.ID = originalEventId;
        ok = cli.updateEvent(event);
    } else {
        ok = cli.addEvent(event, uid, &eventID);
    }

    if (!ok) {
        QMessageBox::warning(this, "Ошибка",
                             QString("Ошибка %1 события").arg(isEditMode ? "изменения" : "создания"));
        ui->create_event->setEnabled(true);
    } else {
        emit eventUpdated();
        if (notifyIndex > 0) {
            ok = cli.scheduleNotification(eventID, event.notificationTime);
        }
        this->close();
    }
}

QDateTime Event::calculateNotificationTime(const QDateTime& eventStart, int notifyIndex)
{
    switch(notifyIndex) {
    case 1: return eventStart.addSecs(-15 * 60);
    case 2: return eventStart.addSecs(-3600);
    case 3: return eventStart.addDays(-1);
    case 4: return eventStart.addDays(-7);
    default: return QDateTime();
    }
}


void Event::fillFormWithEventData(const EventData& data)
{
    ui->nameEdit->setText(data.name);
    ui->descEdit->setText(data.desc);
    ui->typeCombo->setCurrentText(data.type);
    ui->priorityCombo->setCurrentIndex(data.prior - 1);

    ui->startDateEdit->setDate(data.start.date());
    ui->startTimeEdit->setTime(data.start.time());
    ui->endDateEdit->setDate(data.end.date());
    ui->endTimeEdit->setTime(data.end.time());

    // QStringList participantsList;
    // for (const Participant& participant : data.parts) {
    //     participantsList.append(participant.uid);
    // }
    // ui->parts->setPlainText(participantsList.join(" "));

    originalEventId = data.ID;
    isEditMode = true;

}

